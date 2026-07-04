package middleware

import (
	"net/url"
	"slices"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	cfg := config.Get().CORS

	corsCfg := cors.Config{
		AllowMethods:     cfg.Methods,
		AllowHeaders:     cfg.Headers,
		MaxAge:           time.Duration(cfg.MaxAge) * time.Second,
		AllowCredentials: true,
	}

	// 允许所有域名直接设置，否则通过 AllowOriginFunc 匹配主机名（忽略端口，忽略协议）
	if slices.Contains(cfg.Origins, "*") {
		corsCfg.AllowAllOrigins = true
	} else {
		corsCfg.AllowOriginFunc = func(origin string) bool {
			u, err := url.Parse(origin)
			if err != nil {
				return false
			}
			return slices.Contains(cfg.Origins, u.Hostname())
		}
	}

	return cors.New(corsCfg)
}
