package admin

import (
	"errors"
	"fmt"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/captcha"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/permission"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/token"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	repoAdmin "github.com/ai-go-hub/ai-go-admin/internal/repository/admin"
	"github.com/ai-go-hub/ai-go-admin/internal/service"
	"github.com/ai-go-hub/ai-go-admin/pkg/tree"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminService 管理员服务，嵌入通用服务接口并扩展自定义方法
type AdminService struct {
	service.IService[model.Admin]
	repo *repoAdmin.AdminRepository
}

// NewAdminService 创建管理员服务实例
func NewAdminService(repo *repoAdmin.AdminRepository) *AdminService {
	return &AdminService{
		IService: service.NewService(repo),
		repo:     repo,
	}
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string               `json:"username" binding:"required"`
	Password string               `json:"password" binding:"required"`
	Remember bool                 `json:"remember"`
	Captcha  captcha.ClickRequest `json:"captcha"`
}

// LoginResponse 登录响应数据
type LoginResponse struct {
	model.Admin
	Token string `json:"token,omitempty"`
}

// InitResponse 后台初始化响应数据
type InitResponse struct {
	Admin      *dto.AdminSession `json:"admin"`
	Super      bool              `json:"super"`
	SiteConfig map[string]string `json:"site_config"`
	Rules      []map[string]any  `json:"rules"`
}

// Login 管理员登录
func (s *AdminService) Login(c *gin.Context, req *LoginRequest) (*LoginResponse, error) {
	if config.Get().Captcha.Switches.AdminLogin {
		if ok, err := captcha.Check(req.Captcha, true); !ok {
			return nil, fmt.Errorf("验证码错误：%w", err)
		}
	}

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

	// 计算过期时间
	expireCfg := config.Get().Token.Expire
	expiredAt := time.Now().Add(time.Duration(expireCfg.Admin) * time.Hour)
	if req.Remember {
		expiredAt = time.Now().Add(time.Duration(expireCfg.AdminRemember) * time.Hour)
	}

	// 创建令牌
	tk := &model.Token{
		Token:     tokenStr,
		Type:      token.TypeAdminLogin,
		UserID:    admin.ID,
		CreatedAt: time.Now(),
		ExpiredAt: expiredAt,
	}
	if err := token.Manager().Create(c.Request.Context(), tk); err != nil {
		return nil, errors.New("保存令牌失败")
	}

	return &LoginResponse{
		Admin: *admin,
		Token: tokenStr,
	}, nil
}

// Logout 注销当前管理员令牌
func (s *AdminService) Logout(c *gin.Context, tokenStr string) error {
	return token.Manager().Delete(c.Request.Context(), tokenStr)
}

// Init 后台初始化数据聚合
func (s *AdminService) Init(c *gin.Context, adminSession *dto.AdminSession) (*InitResponse, error) {
	ctx := c.Request.Context()

	// 1. 站点配置
	configSiteNames := []string{"name", "record_number", "version"}
	siteConfig := make(map[string]string, len(configSiteNames)+1)

	configs, err := gorm.G[model.Config](database.DB()).
		Where("name IN ?", configSiteNames).
		Find(ctx)
	if err != nil {
		return nil, err
	}
	for _, cfg := range configs {
		siteConfig[cfg.Name] = cfg.Value
	}
	siteConfig["timezone"] = config.Get().App.Timezone

	// 2. 当前管理员拥有的权限规则（菜单）列表
	perm := permission.New()
	rules, err := perm.GetRules(ctx, adminSession.ID)
	if err != nil {
		return nil, err
	}

	// 将菜单规则列表转为树状结构
	// model.AdminRule 转 map[string]any，此方法调用频率较高，使用性能更高的硬编码
	ruleData := make([]map[string]any, len(rules))
	for i, r := range rules {
		pid := uint(0)
		if r.Pid != nil {
			pid = *r.Pid
		}
		ruleData[i] = map[string]any{
			"id":         r.ID,
			"pid":        pid,
			"type":       r.Type,
			"title":      r.Title,
			"name":       r.Name,
			"path":       r.Path,
			"icon":       r.Icon,
			"open_type":  r.OpenType,
			"url":        r.URL,
			"component":  r.Component,
			"keepalive":  r.Keepalive,
			"extend":     r.Extend,
			"remark":     r.Remark,
			"weigh":      r.Weigh,
			"status":     r.Status,
			"updated_at": r.UpdatedAt,
			"created_at": r.CreatedAt,
		}
	}
	ruleTree := tree.Build(ruleData, "id", "pid", "children")

	// 3. 是否为超级管理员
	super, err := perm.IsSuperAdmin(ctx, adminSession.ID)
	if err != nil {
		return nil, err
	}

	return &InitResponse{
		Admin:      adminSession,
		Super:      super,
		SiteConfig: siteConfig,
		Rules:      ruleTree,
	}, nil
}
