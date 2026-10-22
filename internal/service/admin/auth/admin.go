package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/permission"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/bindx"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	repoAdmin "github.com/ai-go-hub/ai-go-admin/internal/repository/admin"
	"github.com/ai-go-hub/ai-go-admin/internal/service"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthAdminService 管理员账号管理服务
type AuthAdminService struct {
	service.IService[model.Admin]
	repo *repoAdmin.AdminRepository
}

// NewAuthAdminService 创建管理员账号管理服务实例
func NewAuthAdminService(repo *repoAdmin.AdminRepository) *AuthAdminService {
	return &AuthAdminService{
		IService: service.NewService(repo),
		repo:     repo,
	}
}

// Create 覆写通用创建方法: 校验用户名唯一 + 加密密码 + 校验分组权限
func (s *AuthAdminService) Create(ctx context.Context, tri *bindx.Tri[model.Admin], opts service.Options) error {
	if tri.Model.Username == "" {
		return errors.New("用户名不能为空")
	}
	if _, err := s.repo.FindByUsername(ctx, tri.Model.Username); err == nil {
		return errors.New("用户名已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if tri.Model.Password == "" {
		return errors.New("密码不能为空")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(tri.Model.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}
	tri.Model.Password = string(hashed)

	// 校验分组权限
	if err := s.validateGroupAccesses(ctx, tri.Model.AdminGroupAccesses, opts); err != nil {
		return err
	}

	return s.IService.Create(ctx, tri, opts)
}

// Update 覆写通用更新方法: 用户名唯一校验 + 密码为空时移除字段不更新，非空则加密 + 分组关联替换
func (s *AuthAdminService) Update(ctx context.Context, tri *bindx.Tri[model.Admin], opts service.Options) error {
	pk, err := strconv.ParseUint(opts.PrimaryKeyValue, 10, 32)
	if err != nil {
		return errors.New("主键值错误")
	}

	// 校验目标管理员是否在操作者可见范围内，防止低权账号操作更高权限账号
	if err := s.ensureTargetsVisible(ctx, []uint{uint(pk)}, opts); err != nil {
		return err
	}

	// 禁止通过管理员更新接口禁用自己的账号
	session, err := sessionFrom(opts)
	if err != nil {
		return err
	}
	if uint(pk) == session.ID && tri.Map["status"] == "disable" {
		return errors.New("不能禁用自己的账号")
	}

	// 用户名唯一校验
	if tri.Model.Username == "" {
		return errors.New("用户名不能为空")
	}
	userNameExist, err := s.repo.FindByUsername(ctx, tri.Model.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if userNameExist != nil && strconv.FormatUint(uint64(userNameExist.ID), 10) != opts.PrimaryKeyValue {
		return errors.New("用户名已存在")
	}

	// 密码为空时移除字段，确保不覆盖旧密码；非空则加密后写入
	// model.Admin.Password 为 json:"-"，但控制器已用 DTO 绑定并经 copier 填入 Model.Password
	if tri.Model.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(tri.Model.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("密码加密失败: %w", err)
		}
		tri.Map["password"] = string(hashed)
	} else {
		delete(tri.Map, "password")
	}

	// 分组信息提取
	accesses := tri.Model.AdminGroupAccesses
	delete(tri.Map, "admin_group_accesses")

	// 校验分组权限
	if err := s.validateGroupAccesses(ctx, accesses, opts); err != nil {
		return err
	}

	if err := s.IService.Update(ctx, tri, opts); err != nil {
		return err
	}

	// 替换分组关联
	if accesses != nil {
		return s.repo.ReplaceGroupAccesses(ctx, uint(pk), accesses)
	}

	return nil
}

// Delete 覆写通用删除方法: 同时删除关联的分组数据
func (s *AuthAdminService) Delete(ctx context.Context, opts service.Options) error {
	if len(opts.PrimaryKeyValues) == 0 {
		return errors.New("请选择要删除的记录")
	}

	uids := make([]uint, 0, len(opts.PrimaryKeyValues))
	for _, pk := range opts.PrimaryKeyValues {
		uid, err := strconv.ParseUint(pk, 10, 32)
		if err != nil {
			return errors.New("主键值错误")
		}
		uids = append(uids, uint(uid))
	}

	session, err := sessionFrom(opts)
	if err != nil {
		return err
	}

	// 禁止删除自己的账号
	if slices.Contains(uids, session.ID) {
		return errors.New("不能删除自己的账号")
	}

	// 校验目标管理员是否在操作者可见范围内，防止低权账号删除更高权限账号
	if err := s.ensureTargetsVisible(ctx, uids, opts); err != nil {
		return err
	}

	// 删除分组关联
	for _, uid := range uids {
		if err := s.repo.DeleteGroupAccesses(ctx, uid); err != nil {
			return err
		}
	}

	return s.IService.Delete(ctx, opts)
}

// ensureTargetsVisible 校验目标管理员的分组规则是否全部在操作者可见范围内
func (s *AuthAdminService) ensureTargetsVisible(ctx context.Context, uids []uint, opts service.Options) error {
	session, err := sessionFrom(opts)
	if err != nil {
		return err
	}

	perm := permission.New()
	super, err := perm.IsSuperAdmin(ctx, session.ID)
	if err != nil {
		return err
	}
	if super {
		return nil
	}

	visible, err := perm.GetVisibleGroupIDs(ctx, session.ID, nil)
	if err != nil {
		return err
	}
	visibleSet := make(map[uint]struct{}, len(visible))
	for _, id := range visible {
		visibleSet[id] = struct{}{}
	}

	targetGroupIDs, err := s.repo.GroupIDsByUIDs(ctx, uids)
	if err != nil {
		return err
	}
	for _, gid := range targetGroupIDs {
		if _, ok := visibleSet[gid]; !ok {
			return errors.New("无权操作该管理员")
		}
	}
	return nil
}

// List 覆写通用查询方法: 只读取自己管理员有权管理的数据
func (s *AuthAdminService) List(ctx context.Context, opts service.Options) ([]model.Admin, error) {
	session, err := sessionFrom(opts)
	if err != nil {
		return nil, err
	}

	perm := permission.New()
	super, err := perm.IsSuperAdmin(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	if !super {
		wheres, err := s.visibilityWheres(ctx, session.ID)
		if err != nil {
			return nil, err
		}
		opts.Wheres = append(opts.Wheres, wheres...)
	}

	return s.IService.List(ctx, opts)
}

// Count 覆写通用统计方法: 只读取自己管理员有权管理的数据
func (s *AuthAdminService) Count(ctx context.Context, opts service.Options) (int64, error) {
	session, err := sessionFrom(opts)
	if err != nil {
		return 0, err
	}

	perm := permission.New()
	super, err := perm.IsSuperAdmin(ctx, session.ID)
	if err != nil {
		return 0, err
	}

	if !super {
		wheres, err := s.visibilityWheres(ctx, session.ID)
		if err != nil {
			return 0, err
		}
		opts.Wheres = append(opts.Wheres, wheres...)
	}

	return s.IService.Count(ctx, opts)
}

// visibilityWheres 返回过滤「不可见管理员」的查询条件
func (s *AuthAdminService) visibilityWheres(ctx context.Context, adminID uint) ([]service.WhereGroup, error) {
	visibleIDs, err := permission.New().GetVisibleGroupIDs(ctx, adminID, nil)
	if err != nil {
		return nil, err
	}

	excludedUIDs, err := s.repo.UIDsByInvisibleGroups(ctx, visibleIDs)
	if err != nil {
		return nil, err
	}
	if len(excludedUIDs) == 0 {
		return nil, nil
	}

	return []service.WhereGroup{{
		Wheres: []service.Where{{
			Field:    "id",
			Operator: "NOT IN",
			Value:    excludedUIDs,
		}},
	}}, nil
}

// validateGroupAccesses 校验分组权限，待分配分组必须全部在操作者可见范围内
func (s *AuthAdminService) validateGroupAccesses(ctx context.Context, accesses []model.AdminGroupAccess, opts service.Options) error {
	if len(accesses) == 0 {
		return nil
	}

	session, err := sessionFrom(opts)
	if err != nil {
		return err
	}

	perm := permission.New()
	super, err := perm.IsSuperAdmin(ctx, session.ID)
	if err != nil {
		return err
	}
	if super {
		return nil
	}

	visible, err := perm.GetVisibleGroupIDs(ctx, session.ID, nil)
	if err != nil {
		return err
	}
	visibleSet := make(map[uint]struct{}, len(visible))
	for _, id := range visible {
		visibleSet[id] = struct{}{}
	}

	for _, acc := range accesses {
		if _, ok := visibleSet[acc.GroupID]; !ok {
			return errors.New("无权分配该分组")
		}
	}
	return nil
}
