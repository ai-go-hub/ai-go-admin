package httpx

import (
	"net"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/gin-gonic/gin"
)

// Scheme 返回当前请求的协议（http / https）
// 直连时依据 TLS 判断；仅当配置了可信代理（server.trusted_proxies）时才依据 X-Forwarded-Proto 头判断
func Scheme(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	// 仅在可信代理环境下信任 X-Forwarded-Proto
	if len(config.Get().Server.TrustedProxies) > 0 {
		if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
			if i := strings.IndexByte(proto, ','); i >= 0 {
				proto = proto[:i]
			}
			scheme = strings.ToLower(strings.TrimSpace(proto))
		}
	}
	return scheme
}

// Port 返回当前请求的端口
// Host 头带端口时直接取；否则按协议返回默认端口（http=80, https=443）
func Port(c *gin.Context) string {
	_, port, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		if Scheme(c) == "https" {
			return "443"
		}
		return "80"
	}
	return port
}

// BaseURL 返回当前站点的根 URL（scheme://host），host 含端口时一并带上
func BaseURL(c *gin.Context) string {
	return Scheme(c) + "://" + c.Request.Host
}

// ClientIP 返回客户端 IP，基于 gin 的 ClientIP（遵循 SetTrustedProxies 配置）
func ClientIP(c *gin.Context) string {
	return c.ClientIP()
}

// ContentType 返回 Content-Type 头
func ContentType(c *gin.Context) string {
	return c.ContentType()
}

// Referer 返回 Referer 头
func Referer(c *gin.Context) string {
	return c.GetHeader("Referer")
}

// UserAgent 返回 User-Agent 头
func UserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}
