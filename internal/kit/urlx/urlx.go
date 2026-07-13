package urlx

import (
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"

	"github.com/gin-gonic/gin"
)

// FullURL 返回静态资源的完整访问 URL
func FullURL(c *gin.Context, resource string) string {
	if resource == "" {
		return ""
	}

	// 已是完整 URL 或 base64 资源，原样返回
	lower := strings.ToLower(resource)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "data:") {
		return resource
	}

	cdn := config.Get().CDN

	// 选择前缀: CDN URL 优先，否则使用当前请求域名
	prefix := cdn.URL
	if prefix == "" {
		prefix = httpx.BaseURL(c)
	}

	// 前缀与路径之间保证恰好一个 /
	u := strings.TrimRight(prefix, "/") + "/" + strings.TrimLeft(resource, "/")

	// 拼接 CDN URL 参数（仅在使用 CDN 时）
	if cdn.URL != "" && cdn.URLParams != "" {
		sep := "?"
		if strings.Contains(u, "?") {
			sep = "&"
		}
		u += sep + cdn.URLParams
	}
	return u
}
