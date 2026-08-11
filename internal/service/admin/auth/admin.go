package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/permission"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/bindx"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	repoAdmin "github.com/ai-go-hub/ai-go-admin/internal/repository/admin"
	repoAuth "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/auth"
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
	groupRepo *repoAuth.AdminGroupRepository
}

// NewAuthAdminService 创建管理员账号管理服务实例
func NewAuthAdminService(repo *repoAdmin.AdminRepository, groupRepo *repoAuth.AdminGroupRepository) *AuthAdminService {
	return &AuthAdminService{
		IService:  service.NewService(repo),
		repo:      repo,
		groupRepo: groupRepo,
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
