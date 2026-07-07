package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/token"
	repoAdmin "github.com/ai-go-hub/ai-go-admin/internal/repository/admin"
	"github.com/ai-go-hub/ai-go-admin/internal/response"

	"github.com/gin-gonic/gin"
)

const (
	// 上下文中存储管理员会话信息的键
	CtxAdminKey = "admin"
)

// 标识常量
const (
	FlagNeedLogin = "need_login" // 需要登录
	FlagLoggedIn  = "logged_in"  // 已经登录
)

// AdminAuth 管理员认证中间件，未登录时阻断请求
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, msg := extractAdmin(c)
		if session != nil {
			c.Set(CtxAdminKey, session)
			c.Next()
		} else {
			response.Fail(c,
				response.WithMessage(msg),
				response.WithCode(http.StatusUnauthorized),
				response.WithData(gin.H{"type": FlagNeedLogin}),
			)
			c.Abort()
		}
	}
}

// AdminAuthOptional 可选管理员认证中间件，有 token 则注入管理员信息，没有也放行
func AdminAuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		if session, _ := extractAdmin(c); session != nil {
			c.Set(CtxAdminKey, session)
		}
		c.Next()
	}
}

// GetAdmin 从上下文中取出管理员信息，未登录时返回 nil
func GetAdmin(c *gin.Context) *dto.AdminSession {
	session, ok := c.Get(CtxAdminKey)
	if !ok {
		return nil
	}
	return session.(*dto.AdminSession)
}

// GetToken 从上下文中取出当前令牌字符串，未登录时返回空
func GetToken(c *gin.Context) string {
	session := GetAdmin(c)
	if session == nil {
		return ""
	}
	return session.Token
}

// extractAdmin 提取并验证 token，返回 (会话信息, 错误消息)
func extractAdmin(c *gin.Context) (*dto.AdminSession, string) {
	// 提取请求 token
	authHeader := c.GetHeader("authorization")

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, "请先登录"
	}

	// 校验令牌
	tk, err := token.Manager().Get(c.Request.Context(), parts[1])
	if err != nil || tk == nil || time.Now().After(tk.ExpiredAt) {
		return nil, "身份认证令牌失效，请重新登录"
	}

	// 仅允许管理员登录类型的令牌
	if tk.Type != token.TypeAdminLogin {
		return nil, "请重新登录"
	}

	// 查询管理员信息
	adminRepo := repoAdmin.NewAdminRepository()
	admin, err := adminRepo.Get(c, tk.UserID)
	if err != nil {
		return nil, "请重新登录"
	}

	// 检查账号状态
	if admin.Status == "disable" {
		return nil, "账号已被禁用"
	}

	return &dto.AdminSession{Admin: admin, Token: parts[1]}, ""
}
