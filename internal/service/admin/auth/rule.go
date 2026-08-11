package auth

import (
	"context"
	"errors"
	"strconv"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/permission"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/bindx"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"
	repoAuth "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/auth"
	"github.com/ai-go-hub/ai-go-admin/internal/service"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"

	"gorm.io/gorm"
)

// AuthAdminRuleExtension 规则列表扩展参数
type AuthAdminRuleExtension struct {
	AdminSession *dto.AdminSession
}

// AuthAdminRuleService 菜单和权限规则管理服务
type AuthAdminRuleService struct {
	service.IService[model.AdminRule]
	repo *repoAuth.AdminRuleRepository
}

// NewAuthAdminRuleService 创建菜单和权限规则管理服务实例
func NewAuthAdminRuleService(repo *repoAuth.AdminRuleRepository) *AuthAdminRuleService {
	return &AuthAdminRuleService{
		IService: service.NewService(repo),
		repo:     repo,
	}
}

// Create 覆写通用创建方法: 增加服务层校验
func (s *AuthAdminRuleService) Create(ctx context.Context, tri *bindx.Tri[model.AdminRule], opts service.Options) error {
	if err := s.validateRule(ctx, &tri.Model, opts.PrimaryKeyValue); err != nil {
		return err
	}
	return s.IService.Create(ctx, tri, opts)
}

// Update 覆写通用更新方法: 增加服务层校验
func (s *AuthAdminRuleService) Update(ctx context.Context, tri *bindx.Tri[model.AdminRule], opts service.Options) error {
	if err := s.validateRule(ctx, &tri.Model, opts.PrimaryKeyValue); err != nil {
		return err
	}
	return s.IService.Update(ctx, tri, opts)
}

// Delete 覆写通用删除方法: 校验被删规则没有游离的子级
func (s *AuthAdminRuleService) Delete(ctx context.Context, opts service.Options) error {
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
	childIDs, err := s.repo.ChildIDsByPids(ctx, pks)
	if err != nil {
		return err
	}

	// 遍历子级，若子级自身不在待删集合内则视为游离节点
	for _, id := range childIDs {
		if _, ok := pkSet[id]; !ok {
			return errors.New("请先删除子级规则，或使用批量删除")
		}
	}

	return s.IService.Delete(ctx, opts)
}

// List 覆写通用查询全部记录方法: 根据管理员权限过滤
func (s *AuthAdminRuleService) List(ctx context.Context, opts service.Options) ([]model.AdminRule, error) {
	// 从控制器传来的 `管理员信息` 扩展数据
	extension, ok := opts.Extension.(*AuthAdminRuleExtension)
	if !ok || extension.AdminSession == nil {
		return nil, errors.New("参数错误，缺少 AdminSession 扩展数据")
	}

	perm := permission.New()
	super, err := perm.IsSuperAdmin(ctx, extension.AdminSession.ID)
	if err != nil {
		return nil, err
	}

	// 无翻页，设定无限 limit
	opts.Limit = 999999

	// 来自选择器则只查 dir、menu，不查 node
	if opts.Selector {
		opts.Wheres = append(opts.Wheres, service.WhereGroup{
			Wheres: []service.Where{{
				Field:    "type",
				Operator: "IN",
				Value:    []string{"dir", "menu"},
			}},
		})
	}

	// 非超管，读取当前管理员拥有的权限规则 IDs
	if !super {
		ruleIDs, err := perm.GetRuleIds(ctx, extension.AdminSession.ID, nil)
		if err != nil {
			return nil, err
		}

		if len(ruleIDs) == 0 {
			return nil, nil
		}

		// 添加 IN IDs 的 where 以确保只读取到当前管理员拥有的权限规则
		opts.Wheres = append(opts.Wheres, service.WhereGroup{
			Wheres: []service.Where{{
				Field:    "id",
				Operator: "IN",
				Value:    ruleIDs,
			}},
		})
	}

	rules, err := s.repo.List(ctx, s.BuildRepoOpts(opts))
	return rules, err
}

// Count 覆写通用统计方法: 与 List 使用相同的权限过滤条件
func (s *AuthAdminRuleService) Count(ctx context.Context, opts service.Options) (int64, error) {
	if opts.Selector {
		return 0, nil
	}

	extension, ok := opts.Extension.(*AuthAdminRuleExtension)
	if !ok || extension.AdminSession == nil {
		return 0, errors.New("参数错误，缺少 AdminSession 扩展数据")
	}

	perm := permission.New()
	rules, err := perm.GetRules(ctx, extension.AdminSession.ID, nil)
	if err != nil {
		return 0, err
	}

	return int64(len(rules)), nil
}

// validateRule 校验规则字段、名称与上级规则
func (s *AuthAdminRuleService) validateRule(ctx context.Context, entity *model.AdminRule, pk string) error {
	if entity.Type == "menu" && util.FromPtr(entity.OpenType) == "tab" && (util.PtrIsZero(entity.Path) || util.PtrIsZero(entity.Component)) {
		return errors.New("规则类型为菜单时，菜单路由路径和菜单组件路径不能为空")
	}

	// 名称唯一校验（排除自身）
	userNameExist, err := s.repo.FindByName(ctx, entity.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if userNameExist != nil && strconv.FormatUint(uint64(userNameExist.ID), 10) != pk {
		return errors.New("规则名称已存在")
	}

	// 非空路径唯一校验（排除自身）
	if util.PtrNotZero(entity.Path) {
		pathExist, err := s.repo.FindByPath(ctx, *entity.Path)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if pathExist != nil && strconv.FormatUint(uint64(pathExist.ID), 10) != pk {
			return errors.New("菜单路由路径已存在")
		}
	}

	// 上级规则校验: 不能将自身设为自己的上级
	if util.PtrNotZero(entity.Pid) {
		strPid := strconv.FormatUint(uint64(*entity.Pid), 10)
		if pk != "" && strPid == pk {
			return errors.New("上级规则不能是自身")
		}
		if _, err := s.repo.Get(ctx, repository.Options{PrimaryKeyValue: strPid}); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("上级规则不存在")
			}
			return err
		}
	}

	return nil
}
