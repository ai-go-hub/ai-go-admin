package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"

	"github.com/didip/tollbooth/v8"
	"github.com/didip/tollbooth/v8/limiter"
	"github.com/gin-gonic/gin"
)

// 全局节流 limiter 单例
var (
	throttleOnce sync.Once
	throttleLmt  *limiter.Limiter
)

// buildThrottleLimiter 根据配置构建限流器，并启动过期令牌桶清理 goroutine
func buildThrottleLimiter(cfg config.ThrottleConfig) *limiter.Limiter {
	lmt := tollbooth.NewLimiter(cfg.Max, &limiter.ExpirableOptions{
		// 总过期时间
		DefaultExpirationTTL: time.Duration(cfg.ExpirationTTL) * time.Second,
	}).
		// tollbooth v8 仅支持单一 IPLookup，固定从 RemoteAddr 取值，
		// 真实客户端 IP 由 httpx.ClientIP 解析后写入
		SetIPLookup(limiter.IPLookup{Name: "RemoteAddr"}).
		SetBurst(cfg.Burst).
		SetIgnoreURL(cfg.IgnoreURL).
		SetMethods(cfg.Methods).
		SetHeaders(cfg.Headers).
		SetMessage(cfg.Message).
		SetMessageContentType(cfg.MessageContentType).
		SetStatusCode(cfg.StatusCode)

	// 底层 expirable-cache 无自动清理机制，过期令牌桶会滞留内存，按库建议以 TTL/2 周期清理
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.ExpirationTTL) * time.Second / 2)
		defer ticker.Stop()
		for range ticker.C {
			lmt.DeleteExpiredTokenBuckets()
		}
	}()

	return lmt
}

// Throttle 全局接口节流中间件
func Throttle() gin.HandlerFunc {
	cfg := config.Get().Throttle

	// 未启用节流时直接放行
	if !cfg.Enabled {
		return func(c *gin.Context) { c.Next() }
	}

	// 单例构建 limiter 与清理 goroutine
	throttleOnce.Do(func() {
		throttleLmt = buildThrottleLimiter(cfg)
	})

	lmt := throttleLmt
	isJSON := strings.Contains(lmt.GetMessageContentType(), "json")

	return func(c *gin.Context) {
		ip := httpx.ClientIP(c)

		// 无法解析出客户端 IP 时直接拒绝（fail-closed，防止伪造请求绕过限流）
		// 真实服务里 RemoteAddr 恒为 host:port，所以此处仅为防御性兜底
		if ip == "" {
			c.JSON(cfg.StatusCode, httpx.Response{
				Code: 1, Message: "无法解析 IP", Time: time.Now().Unix(),
			})
			c.Abort()
			return
		}

		// 浅拷贝请求并把解析出的 IP 写入 RemoteAddr
		r := new(http.Request)
		*r = *c.Request
		r.RemoteAddr = net.JoinHostPort(ip, "0")

		// 跳过限流
		if tollbooth.ShouldSkipLimiter(lmt, r) {
			c.Next()
			return
		}

		// 不走 LimitByRequest，避免给每个响应都注入 X-Rate-Limit-* / RateLimit-* 头
		for _, keys := range tollbooth.BuildKeys(lmt, r) {
			httpError := tollbooth.LimitByKeys(lmt, keys)
			if httpError != nil {
				// 按 message_content_type 配置决定响应格式: json 走项目统一响应结构，其余走原始文本
				if isJSON {
					c.JSON(httpError.StatusCode, httpx.Response{
						Code:    1,
						Message: httpError.Message,
						Time:    time.Now().Unix(),
					})
				} else {
					c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
				}
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
