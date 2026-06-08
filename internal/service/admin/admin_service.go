package admin

import (
	"errors"

	"ai-go-mall/internal/model"
	repoAdmin "ai-go-mall/internal/repository/admin"
	"ai-go-mall/internal/service"

	"github.com/gin-gonic/gin"
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
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
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

	return &LoginResponse{
		Admin: *admin,
		// TODO: 生成 Token（待实现认证模块）
		Token: "",
	}, nil
}
