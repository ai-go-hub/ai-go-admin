package dto

import (
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/captcha"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/pkg/copierx"

	"github.com/jinzhu/copier"
)

// Admin 管理员信息
type Admin struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	LastLoginAt string `json:"last_login_at"`
	LastLoginIP string `json:"last_login_ip"`
	Bio         string `json:"bio"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

// AdminSession 管理员会话信息，由认证中间件注入到请求上下文
// 管理员信息任旧是 *model.Admin 并非 dto.Admin
type AdminSession struct {
	*model.Admin
	Token string `json:"token"`
}

// NewAdmin 将 *model.Admin 转换为 *dto.Admin，
// time.Time 类型的字段会被转换为 string，时间格式为 time.DateTime（yyyy-mm-dd h:i:s）
func NewAdmin(m *model.Admin) (*Admin, error) {
	d := &Admin{}
	if err := copier.CopyWithOption(d, m, copier.Option{
		Converters: []copier.TypeConverter{copierx.Time(time.DateTime)},
	}); err != nil {
		return nil, err
	}
	return d, nil
}

// LoginRequest 管理员登录请求参数
type LoginRequest struct {
	Username string               `json:"username" form:"username" binding:"required"`
	Password string               `json:"password" form:"password" binding:"required"`
	Remember bool                 `json:"remember" form:"remember"`
	Captcha  captcha.ClickRequest `json:"captcha" form:"captcha"`
}

// LoginResponse 管理员登录响应数据
type LoginResponse struct {
	*Admin
	Token string `json:"token,omitempty"`
}

// InitResponse 后台初始化响应数据
type InitResponse struct {
	Admin      *Admin            `json:"admin"`
	Super      bool              `json:"super"`
	SiteConfig map[string]string `json:"site_config"`
	Rules      []map[string]any  `json:"rules"`
}
