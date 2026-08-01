package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/permission"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	repoAdmin "github.com/ai-go-hub/ai-go-admin/internal/repository/admin"
	"github.com/ai-go-hub/ai-go-admin/internal/service"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthAdminExtension 管理员操作扩展参数
type AuthAdminExtension struct {
	AdminSession *dto.AdminSession
}

// AuthAdminService 管理员账号管理服务
type AuthAdminService struct {
	service.IService[model.Admin]
	repo      *repoAdmin.AdminRepository
	groupRepo *repoAdmin.AdminGroupRepository
}

// NewAuthAdminService 创建管理员账号管理服务实例
func NewAuthAdminService(repo *repoAdmin.AdminRepository, groupRepo *repoAdmin.AdminGroupRepository) *AuthAdminService {
	return &AuthAdminService{
		IService:  service.NewService(repo),
		repo:      repo,
		groupRepo: groupRepo,
	}
}

// Create 覆写通用创建方法: 校验用户名唯一 + 加密密码 + 校验分组权限
func (s *AuthAdminService) Create(ctx context.Context, entity *model.Admin, opts service.Options) error {
	if entity.Username == "" {
		return errors.New("用户名不能为空")
	}
	if _, err := s.repo.FindByUsername(ctx, entity.Username); err == nil {
		return errors.New("用户名已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if entity.Password == "" {
		return errors.New("密码不能为空")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(entity.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}
	entity.Password = string(hashed)

	// 校验分组权限
	if err := s.validateGroupAccesses(ctx, entity.AdminGroupAccesses, opts); err != nil {
		return err
	}

	return s.IService.Create(ctx, entity, opts)
}

// Update 覆写通用更新方法: 用户名唯一校验 + 密码为空时跳过更新，非空则加密 + 校验分组权限
// 依赖 GORM struct 更新时会跳过零值字段的特性，因此密码为空不会覆盖旧值
func (s *AuthAdminService) Update(ctx context.Context, entity *model.Admin, opts service.Options) error {
	if entity.Username == "" {
		return errors.New("用户名不能为空")
	}
	userNameExist, err := s.repo.FindByUsername(ctx, entity.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if userNameExist != nil && strconv.FormatUint(uint64(userNameExist.ID), 10) != opts.PrimaryKeyValue {
		return errors.New("用户名已存在")
	}

	if entity.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(entity.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("密码加密失败: %w", err)
		}
		entity.Password = string(hashed)
	}

	pk, err := strconv.ParseUint(opts.PrimaryKeyValue, 10, 32)
	if err != nil {
		return errors.New("主键值错误")
	}

	// 保存分组关联数据，以便后续先清空再插入
	accesses := entity.AdminGroupAccesses

	// 置空以免基类自动处理关联数据
	entity.AdminGroupAccesses = nil

	// 校验分组权限
	if err := s.validateGroupAccesses(ctx, accesses, opts); err != nil {
		return err
	}

	if err := s.IService.Update(ctx, entity, opts); err != nil {
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

	// 先删除分组关联
	for _, pk := range opts.PrimaryKeyValues {
		uid, err := strconv.ParseUint(pk, 10, 32)
		if err != nil {
			return errors.New("主键值错误")
		}
		if err := s.repo.DeleteGroupAccesses(ctx, uint(uid)); err != nil {
			return err
		}
	}

	return s.IService.Delete(ctx, opts)
}

// validateGroupAccesses 校验分组权限: 操作管理员必须拥有待分配分组的全部权限规则
func (s *AuthAdminService) validateGroupAccesses(ctx context.Context, accesses []model.AdminGroupAccess, opts service.Options) error {
	if len(accesses) == 0 {
		return nil
	}

	ext, ok := opts.Extension.(*AuthAdminExtension)
	if !ok || ext.AdminSession == nil {
		return errors.New("参数错误，缺少 AdminSession 扩展数据")
	}
	session := ext.AdminSession

	perm := permission.New()
	super, err := perm.IsSuperAdmin(ctx, session.ID)
	if err != nil {
		return err
	}
	if super {
		return nil
	}

	// 获取当前管理员的全部规则 ID
	myRuleIDs, err := perm.GetRuleIds(ctx, session.ID, nil)
	if err != nil {
		return err
	}
	myRuleSet := make(map[uint]struct{}, len(myRuleIDs))
	for _, id := range myRuleIDs {
		myRuleSet[id] = struct{}{}
	}

	// 收集待分配的分组 ID
	groupIDs := make([]uint, len(accesses))
	for i, acc := range accesses {
		groupIDs[i] = acc.GroupID
	}

	groups, err := s.groupRepo.FindByIDs(ctx, groupIDs)
	if err != nil {
		return err
	}

	for _, group := range groups {
		if group.Rules == nil || *group.Rules == "" {
			continue
		}
		if *group.Rules == "*" {
			return fmt.Errorf("无权分配超级管理员分组 '%s'", group.Name)
		}
		groupRuleIDs, err := parseRuleIDs(group.Rules)
		if err != nil {
			return err
		}
		for _, rid := range groupRuleIDs {
			if _, ok := myRuleSet[rid]; !ok {
				return fmt.Errorf("无权分配分组 '%s'，您缺少该分组的部分权限", group.Name)
			}
		}
	}

	return nil
}
