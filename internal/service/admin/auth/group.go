package auth

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/permission"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"
	repoAdmin "github.com/ai-go-hub/ai-go-admin/internal/repository/admin"
	"github.com/ai-go-hub/ai-go-admin/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ai-go-hub/ai-go-admin/pkg/util"
)

// AuthAdminGroupExtension 分组列表扩展参数
type AuthAdminGroupExtension struct {
	AdminSession *dto.AdminSession
}

// AuthAdminGroupService 管理员分组管理服务
type AuthAdminGroupService struct {
	service.IService[model.AdminGroup]
	repo     *repoAdmin.AdminGroupRepository
	ruleRepo *repoAdmin.AdminRuleRepository
}

// NewAuthAdminGroupService 创建管理员分组管理服务实例
func NewAuthAdminGroupService(repo *repoAdmin.AdminGroupRepository, ruleRepo *repoAdmin.AdminRuleRepository) *AuthAdminGroupService {
	return &AuthAdminGroupService{
		IService: service.NewService(repo),
		repo:     repo,
		ruleRepo: ruleRepo,
	}
}

// Create 覆写通用创建方法
func (s *AuthAdminGroupService) Create(c *gin.Context, entity *model.AdminGroup, opts service.Options) error {
	// 分组基本信息校验
	if err := s.validateGroup(c, entity, opts); err != nil {
		return err
	}

	// 若选中的规则已覆盖全部规则，折叠为通配符 "*"
	hasAll, err := hasAllRules(c, entity)
	if err != nil {
		return err
	}
	if hasAll {
		entity.Rules = util.ToPtr("*")
	}

	// 权限规则校验
	if err := s.validateRules(c, entity, opts); err != nil {
		return err
	}
	return s.IService.Create(c, entity, opts)
}

// Update 覆写通用更新方法
func (s *AuthAdminGroupService) Update(c *gin.Context, entity *model.AdminGroup, opts service.Options) error {
	// 分组基本信息校验
	if err := s.validateGroup(c, entity, opts); err != nil {
		return err
	}

	// 不允许将分组挂到自身的后代分组下，避免形成环
	if pid := util.FromPtr(entity.Pid); pid != 0 {
		if err := s.ensurePidNotDescendant(c, opts.PrimaryKeyValue, pid); err != nil {
			return err
		}
	}

	// 若选中的规则已覆盖全部规则，折叠为通配符 "*"
	hasAll, err := hasAllRules(c, entity)
	if err != nil {
		return err
	}
	if hasAll {
		entity.Rules = util.ToPtr("*")
	}

	// 权限规则校验
	if err := s.validateRules(c, entity, opts); err != nil {
		return err
	}
	return s.IService.Update(c, entity, opts)
}

// Delete 覆写通用删除方法: 校验被删分组没有游离的子级
func (s *AuthAdminGroupService) Delete(c *gin.Context, opts service.Options) error {
	if len(opts.PrimaryKeyValues) == 0 {
		return errors.New("请选择要删除的记录")
	}

	// 把待删主键转为 uint 集合（与数据库 pid 字段类型一致，避免类型转换开销）
	pks := make([]uint, 0, len(opts.PrimaryKeyValues))
	pkSet := make(map[uint]struct{}, len(opts.PrimaryKeyValues))
	for _, pk := range opts.PrimaryKeyValues {
		id, err := strconv.ParseUint(pk, 10, 32)
		if err != nil {
			return errors.New("主键值错误")
		}
		pks = append(pks, uint(id))
		pkSet[uint(id)] = struct{}{}
	}

	// 查询 pid 的直接子级 id
	childIDs, err := s.repo.ChildIDsByPids(c, pks)
	if err != nil {
		return err
	}

	// 遍历子级，若子级自身不在待删集合内则视为游离子级
	for _, id := range childIDs {
		if _, ok := pkSet[id]; !ok {
			return errors.New("请先删除子级分组，或使用批量删除")
		}
	}

	return s.IService.Delete(c, opts)
}

// List 覆写通用查询全部记录方法: 根据管理员权限过滤（自身所属分组 + 其后代分组）
func (s *AuthAdminGroupService) List(c *gin.Context, opts service.Options) ([]model.AdminGroup, error) {
	extension, ok := opts.Extension.(*AuthAdminGroupExtension)
	if !ok || extension.AdminSession == nil {
		return nil, errors.New("参数错误，缺少 AdminSession 扩展数据")
	}

	// 树状表格无翻页，设定无限 limit
	opts.Limit = 999999

	perm := permission.New()
	super, err := perm.IsSuperAdmin(c.Request.Context(), extension.AdminSession.ID)
	if err != nil {
		return nil, err
	}

	// 非超管，需要限定到读取 `自身所属分组 + 其后代分组` 的规则
	if !super {
		visibleIDs, err := s.visibleGroupIDs(c, extension.AdminSession.ID)
		if err != nil {
			return nil, err
		}
		if len(visibleIDs) == 0 {
			return nil, nil
		}
		opts.Wheres = append(opts.Wheres, service.WhereGroup{
			Wheres: []service.Where{{
				Field:    "id",
				Operator: "IN",
				Value:    visibleIDs,
			}},
		})
	}

	return s.repo.List(c, s.BuildRepoOpts(opts))
}

// Count 覆写通用统计方法: 与 List 使用相同的权限过滤条件
func (s *AuthAdminGroupService) Count(c *gin.Context, opts service.Options) (int64, error) {
	if opts.Selector {
		return 0, nil
	}

	extension, ok := opts.Extension.(*AuthAdminGroupExtension)
	if !ok || extension.AdminSession == nil {
		return 0, errors.New("参数错误，缺少 AdminSession 扩展数据")
	}

	perm := permission.New()
	super, err := perm.IsSuperAdmin(c.Request.Context(), extension.AdminSession.ID)
	if err != nil {
		return 0, err
	}

	if !super {
		visibleIDs, err := s.visibleGroupIDs(c, extension.AdminSession.ID)
		if err != nil {
			return 0, err
		}
		if len(visibleIDs) == 0 {
			return 0, nil
		}
		opts.Wheres = append(opts.Wheres, service.WhereGroup{
			Wheres: []service.Where{{
				Field:    "id",
				Operator: "IN",
				Value:    visibleIDs,
			}},
		})
	}

	return s.IService.Count(c, opts)
}

// StripParentRuleIDs 计算剥离父级 ID 之后的 rules 值（父级选择状态由子级决定）
func (s *AuthAdminGroupService) StripParentRuleIDs(c *gin.Context, entity *model.AdminGroup) (*string, error) {
	ruleIDs, err := parseRuleIDs(entity.Rules)
	if err != nil || len(ruleIDs) == 0 {
		return nil, err
	}

	pids, err := s.ruleRepo.DistinctPidsByIDs(c, ruleIDs)
	if err != nil {
		return nil, err
	}
	if len(pids) == 0 {
		return entity.Rules, nil
	}

	pidSet := make(map[uint]struct{}, len(pids))
	for _, id := range pids {
		pidSet[id] = struct{}{}
	}

	// 剔除作为父级出现的 ID
	kept := make([]string, 0, len(ruleIDs))
	for _, id := range ruleIDs {
		if _, isParent := pidSet[id]; isParent {
			continue
		}
		kept = append(kept, strconv.FormatUint(uint64(id), 10))
	}
	return util.ToPtr(strings.Join(kept, ",")), nil
}

// BuildRulesTitles 为分组列表构建 rules 字段的摘要文本
func (s *AuthAdminGroupService) BuildRulesTitles(c *gin.Context, groups []model.AdminGroup) (map[uint]string, error) {
	result := make(map[uint]string, len(groups))

	for _, g := range groups {
		if util.FromPtr(g.Rules) == "*" {
			result[g.ID] = "超级管理员"
			continue
		}

		ruleIDs, err := parseRuleIDs(g.Rules)
		if err != nil || len(ruleIDs) == 0 {
			result[g.ID] = "无权限"
			continue
		}

		title, err := s.firstMenuTitle(c, ruleIDs)
		if err != nil {
			return nil, err
		}

		if len(ruleIDs) == 1 {
			result[g.ID] = title
		} else {
			result[g.ID] = fmt.Sprintf("%s等 %d 项", title, len(ruleIDs))
		}
	}
	return result, nil
}

// validateGroup 校验分组名称、上级分组
func (s *AuthAdminGroupService) validateGroup(c *gin.Context, entity *model.AdminGroup, opts service.Options) error {
	// 名称唯一校验
	nameExist, err := s.repo.FindByName(c, entity.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if nameExist != nil && strconv.FormatUint(uint64(nameExist.ID), 10) != opts.PrimaryKeyValue {
		return errors.New("组名已存在")
	}

	// 上级分组校验
	if util.FromPtr(entity.Pid) != 0 {
		strPid := strconv.FormatUint(uint64(*entity.Pid), 10)
		if strPid == opts.PrimaryKeyValue {
			return errors.New("上级分组不能是自身")
		}
		if _, err := s.repo.Get(c, repository.Options{PrimaryKeyValue: strPid}); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("上级分组不存在")
			}
			return err
		}
	}

	return nil
}

// validateRules 校验 rules 字段格式与越权授权
//
// 超管级分组仅允许超级管理员设置，非超管所填的规则 ID 必须都在自己拥有的规则集合内
func (s *AuthAdminGroupService) validateRules(c *gin.Context, entity *model.AdminGroup, opts service.Options) error {
	extension, _ := opts.Extension.(*AuthAdminGroupExtension)
	if extension == nil || extension.AdminSession == nil {
		return errors.New("参数错误，缺少 AdminSession 扩展数据")
	}

	perm := permission.New()
	super, err := perm.IsSuperAdmin(c.Request.Context(), extension.AdminSession.ID)
	if err != nil {
		return err
	}

	// 通配符仅超管可用
	if strings.Contains(util.FromPtr(entity.Rules), "*") {
		if !super {
			return errors.New("无权设置超级管理员分组")
		}

		// 是超管直接放行
		return nil
	}

	ruleIDs, err := parseRuleIDs(entity.Rules)
	if err != nil {
		return err
	}
	if len(ruleIDs) == 0 || super {
		return nil
	}

	// 非超管: 分组的 rules 必须是当前操作人拥有的规则集合的子集
	ownedIDs, err := perm.GetRuleIds(c.Request.Context(), extension.AdminSession.ID, nil)
	if err != nil {
		return err
	}
	ownedSet := make(map[uint]struct{}, len(ownedIDs))
	for _, id := range ownedIDs {
		ownedSet[id] = struct{}{}
	}
	for _, id := range ruleIDs {
		if _, ok := ownedSet[id]; !ok {
			return errors.New("不能授予自己没有的权限规则")
		}
	}

	return nil
}

// visibleGroupIDs 返回当前管理员可见的分组 ID
func (s *AuthAdminGroupService) visibleGroupIDs(c *gin.Context, adminID uint) ([]uint, error) {
	perm := permission.New()

	// 一般超管不会再进入此方法，此处加入 super 相关逻辑仅为托底（所以不为超管写单独的短路条件）
	super, err := perm.IsSuperAdmin(c.Request.Context(), adminID)
	if err != nil {
		return nil, err
	}

	// 拿到当前管理员的规则 ID 集合
	ownedIDs, err := perm.GetRuleIds(c.Request.Context(), adminID, nil)
	if err != nil {
		return nil, err
	}
	if len(ownedIDs) == 0 {
		return nil, nil
	}

	var ownedSet map[uint]struct{}
	ownedSet = make(map[uint]struct{}, len(ownedIDs))
	for _, id := range ownedIDs {
		ownedSet[id] = struct{}{}
	}

	all, err := s.repo.FindAll(c)
	if err != nil {
		return nil, err
	}

	// isSubset 判断 ids <= ownedSet（范围）
	isSubset := func(ids []uint) bool {
		for _, id := range ids {
			if _, ok := ownedSet[id]; !ok {
				return false
			}
		}
		return true
	}

	result := make([]uint, 0, len(all))
	for _, g := range all {
		if super {
			result = append(result, g.ID)
			continue
		}

		// 超管级分组: 非超管一律不可见
		if util.FromPtr(g.Rules) == "*" {
			continue
		}

		ruleIDs, err := parseRuleIDs(g.Rules)
		if err != nil {
			// 无法解析视为不可见
			continue
		}

		if isSubset(ruleIDs) {
			result = append(result, g.ID)
		}
	}
	return result, nil
}

// firstMenuTitle 按 ids 传入顺序找出首个 type=menu|dir 节点的 title
func (s *AuthAdminGroupService) firstMenuTitle(c *gin.Context, ids []uint) (string, error) {
	titleMap, err := s.ruleRepo.TitleMapByIDs(c, ids, "menu", "dir")
	if err != nil {
		return "", err
	}
	if title := pickTitleByOrder(ids, titleMap); title != "" {
		return title, nil
	}

	// 没有找到有效的 title，获取 ids 的 pids
	pids, err := s.ruleRepo.DistinctPidsByIDs(c, ids)
	if err != nil {
		return "", err
	}
	if len(pids) == 0 {
		return "", nil
	}

	// 以 pids 再次获取有效的 title
	titleMap, err = s.ruleRepo.TitleMapByIDs(c, pids, "menu", "dir")
	if err != nil {
		return "", err
	}
	return pickTitleByOrder(pids, titleMap), nil
}

// ensurePidNotDescendant 校验新的 pid 不是当前分组的后代分组
func (s *AuthAdminGroupService) ensurePidNotDescendant(c *gin.Context, pk string, pid uint) error {
	all, err := s.repo.FindAll(c)
	if err != nil {
		return err
	}
	selfID, err := strconv.ParseUint(pk, 10, 32)
	if err != nil {
		return errors.New("主键值错误")
	}

	// 收集所有子级分组 ID
	descendants := descendantIDs(all, uint(selfID))
	if slices.Contains(descendants, pid) {
		return errors.New("上级分组不能是当前分组的下级")
	}
	return nil
}

// descendantIDs 返回 rootID 在 all 中的所有子级分组 ID（不包含 rootID 自身）
func descendantIDs(all []model.AdminGroup, rootID uint) []uint {
	// pid -> children ids
	childrenMap := make(map[uint][]uint, len(all))
	for _, g := range all {
		if util.FromPtr(g.Pid) == 0 {
			continue
		}
		childrenMap[*g.Pid] = append(childrenMap[*g.Pid], g.ID)
	}

	result := make([]uint, 0)
	stack := append([]uint{}, childrenMap[rootID]...)

	for len(stack) > 0 {
		id := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, id)
		stack = append(stack, childrenMap[id]...)
	}
	return result
}

// parseRuleIDs 解析 rules 字符串为 ID 切片
func parseRuleIDs(rules *string) ([]uint, error) {
	if rules == nil {
		return nil, nil
	}
	s := strings.TrimSpace(*rules)
	if s == "" || strings.Contains(s, "*") {
		return nil, nil
	}
	ids := make([]uint, 0)
	for part := range strings.SplitSeq(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return nil, errors.New("权限规则 ID 格式错误")
		}
		ids = append(ids, uint(id))
	}
	return ids, nil
}

// hasAllRules 判断 entity.Rules 中的规则 ID 是否已覆盖数据库中的全部规则
// 仅做数据判定，不改变 entity，也不做任何权限校验
func hasAllRules(c *gin.Context, entity *model.AdminGroup) (bool, error) {
	ruleIDs, err := parseRuleIDs(entity.Rules)
	if err != nil || len(ruleIDs) == 0 {
		return false, err
	}

	allIDs, err := permission.New().AllRuleIDs(c.Request.Context(), nil)
	if err != nil {
		return false, err
	}
	if len(allIDs) == 0 || len(ruleIDs) < len(allIDs) {
		return false, nil
	}

	// 使用 set 判断 ruleIDs 是否包含全部 allIDs（去重后逐一比对）
	provided := make(map[uint]struct{}, len(ruleIDs))
	for _, id := range ruleIDs {
		provided[id] = struct{}{}
	}
	for _, id := range allIDs {
		if _, ok := provided[id]; !ok {
			return false, nil
		}
	}

	return true, nil
}

// pickTitleByOrder 按 ids 顺序返回 titleMap 中首个命中的 title
func pickTitleByOrder(ids []uint, titleMap map[uint]string) string {
	for _, id := range ids {
		if t, ok := titleMap[id]; ok {
			return t
		}
	}
	return ""
}
