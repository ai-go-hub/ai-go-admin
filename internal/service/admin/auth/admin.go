package auth

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	repoAdmin "github.com/ai-go-hub/ai-go-admin/internal/repository/admin"
	"github.com/ai-go-hub/ai-go-admin/internal/service"

	"github.com/gin-gonic/gin"
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

// Create 覆写通用创建方法: 校验用户名唯一 + 加密密码
func (s *AuthAdminService) Create(c *gin.Context, entity *model.Admin, opts service.Options) error {
	if entity.Username == "" {
		return errors.New("用户名不能为空")
	}
	if _, err := s.repo.FindByUsername(c, entity.Username); err == nil {
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

	return s.IService.Create(c, entity, opts)
}

// Update 覆写通用更新方法: 用户名唯一校验 + 密码为空时跳过更新，非空则加密
// 依赖 GORM struct 更新时会跳过零值字段的特性，因此密码为空不会覆盖旧值
func (s *AuthAdminService) Update(c *gin.Context, entity *model.Admin, opts service.Options) error {
	if entity.Username == "" {
		return errors.New("用户名不能为空")
	}
	userNameExist, err := s.repo.FindByUsername(c, entity.Username)
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

	return s.IService.Update(c, entity, opts)
}
