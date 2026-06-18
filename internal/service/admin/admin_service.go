package admin

import (
	"errors"
	"time"

	"ai-go-mall/internal/infra/token"
	"ai-go-mall/internal/model"
	repoAdmin "ai-go-mall/internal/repository/admin"
	"ai-go-mall/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service 管理员服务，嵌入通用服务接口并扩展自定义方法
type Service struct {
	service.IService[model.Admin]
	repo *repoAdmin.Repository
}

// NewService 创建管理员服务实例
func NewService(repo *repoAdmin.Repository) *Service {
	return &Service{
		IService: service.NewService(repo),
		repo:     repo,
	}
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remember bool   `json:"remember"`
}

// LoginResponse 登录响应数据
type LoginResponse struct {
	model.Admin
	Token string `json:"token,omitempty"`
}

// Login 管理员登录
func (s *Service) Login(c *gin.Context, req *LoginRequest) (*LoginResponse, error) {
	// 根据用户名查询管理员
	admin, err := s.repo.FindByUsername(c, req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查账号状态
	if admin.Status == "disable" {
		return nil, errors.New("账号已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		// 增加登录失败次数
		_ = s.repo.IncrementLoginFailure(c, admin.ID)
		return nil, errors.New("用户名或密码错误")
	}

	// 更新登录信息
	_ = s.repo.UpdateLoginInfo(c, admin.ID, c.ClientIP())

	// 使用 UUID v7 生成令牌
	tokenStr := uuid.Must(uuid.NewV7()).String()

	// 计算过期时间：记住登录为 30 天，否则为 3 天
	expiredAt := time.Now().Add(3 * 24 * time.Hour)
	if req.Remember {
		expiredAt = time.Now().Add(30 * 24 * time.Hour)
	}

	// 创建令牌
	tk := &model.Token{
		Token:     tokenStr,
		Type:      "admin_login",
		UserID:    admin.ID,
		CreatedAt: time.Now(),
		ExpiredAt: expiredAt,
	}
	if err := token.Instance().Create(c.Request.Context(), tk); err != nil {
		return nil, errors.New("保存令牌失败")
	}

	return &LoginResponse{
		Admin: *admin,
		Token: tokenStr,
	}, nil
}
