package admin

import (
	"context"
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
	"github.com/ai-go-hub/ai-go-admin/pkg/util"

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

// errWrongCredentials 统一的登录失败文案: 用户名不存在、密码错误均返回同一错误，避免用户名枚举
var errWrongCredentials = errors.New("用户名或密码错误")

// 管理员登录安全参数
const (
	// maxLoginFailures 连续登录失败次数阈值，达到后锁定账号
	maxLoginFailures = 10
	// loginLockDuration 连续登录失败的累计窗口，同时也是账号锁定时长
	loginLockDuration = 24 * time.Hour
)

// Login 管理员登录
func (s *AdminService) Login(ctx context.Context, req *dto.LoginRequest, clientIP string) (*dto.LoginResponse, error) {
	if config.Get().Captcha.Switches.AdminLogin {
		if ok, err := captcha.Check(req.Captcha, true); !ok {
			return nil, fmt.Errorf("验证码错误：%w", err)
		}
	}

	// 根据用户名查询管理员
	admin, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		// 用户名不存在时也执行一次等价开销的 bcrypt 运算，抹平与真实密码比对的响应时间差异。
		// cost 与创建密码时一致；bcrypt 仅处理前 72 字节，超长部分截断，
		// 否则 GenerateFromPassword 对超长密码会直接报错跳过运算，时间差异重现
		password := []byte(req.Password)
		if len(password) > 72 {
			password = password[:72]
		}
		_, _ = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		return nil, errWrongCredentials
	}

	// 失败窗口过期判断：距上次登录尝试超过锁定时长则清零失败计数重新累计
	if admin.LoginFailure > 0 && admin.LastLoginAt != nil && time.Since(*admin.LastLoginAt) > loginLockDuration {
		_ = s.repo.ResetLoginFailure(ctx, admin.ID)
		admin.LoginFailure = 0
	}

	// 锁定中直接拒绝：不校验密码、不写库，时间锚点保持冻结，锁定时长固定为 loginLockDuration；
	// 窗口过期后上方已清零计数，自然放行。
	if admin.LoginFailure >= maxLoginFailures {
		return nil, errors.New("请求频繁，请改天再试")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		// 记录失败次数并刷新时间锚点，达到阈值后锚点冻结，账号进入锁定
		_ = s.repo.RecordLoginFailure(ctx, admin.ID, clientIP, maxLoginFailures)
		return nil, errWrongCredentials
	}

	// 检查账号状态
	if admin.Status == "disable" {
		return nil, errors.New("账号已被禁用")
	}

	// 更新登录信息
	_ = s.repo.UpdateLoginInfo(ctx, admin.ID, clientIP)

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
	if err := token.Manager().Create(ctx, tk); err != nil {
		return nil, errors.New("保存令牌失败")
	}

	adminDTO, err := dto.NewAdmin(admin)
	if err != nil {
		return nil, errors.New("数据转 DTO 失败")
	}
	return &dto.LoginResponse{Admin: adminDTO, Token: tokenStr}, nil
}

// Logout 注销当前管理员令牌
func (s *AdminService) Logout(ctx context.Context, tokenStr string) error {
	return token.Manager().Delete(ctx, tokenStr)
}

// Init 后台初始化数据聚合
func (s *AdminService) Init(ctx context.Context, adminSession *dto.AdminSession) (*dto.InitResponse, error) {
	// 1. 站点配置
	configSiteNames := []string{"name", "record_number", "ps_record_number", "version"}
	siteConfig := make(map[string]string, len(configSiteNames)+3)

	configs, err := gorm.G[model.Config](database.DB()).
		Where("name IN ?", configSiteNames).
		Find(ctx)
	if err != nil {
		return nil, err
	}
	for _, cfg := range configs {
		siteConfig[cfg.Name] = util.FromPtr(cfg.Value)
	}
	siteConfig["timezone"] = config.Get().App.Timezone
	siteConfig["cdn_url"] = config.Get().CDN.URL
	siteConfig["cdn_url_params"] = config.Get().CDN.URLParams

	// 2. 当前管理员拥有的权限规则（菜单）列表
	perm := permission.New()
	rules, err := perm.GetRules(ctx, adminSession.ID, new(permission.RuleStatusEnabled))
	if err != nil {
		return nil, err
	}

	// 将菜单规则列表转为树状结构
	ruleData := make([]map[string]any, len(rules))
	for i := range rules {
		ruleData[i] = rules[i].ToMap()
	}
	ruleTree := tree.Build(ruleData, "id", "pid", "children")

	// 3. 是否为超级管理员
	super, err := perm.IsSuperAdmin(ctx, adminSession.ID)
	if err != nil {
		return nil, err
	}

	// 4. 将 AdminSession 内的 *model.Admin 转换为 DTO，转换期间会格式化时间
	adminDTO, err := dto.NewAdmin(adminSession.Admin)
	if err != nil {
		return nil, errors.New("数据转 DTO 失败")
	}

	return &dto.InitResponse{
		Admin:      adminDTO,
		Super:      super,
		SiteConfig: siteConfig,
		Rules:      ruleTree,
	}, nil
}
