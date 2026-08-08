package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/permission"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"
	repoAuth "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/auth"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"

	"github.com/gin-gonic/gin"
)

const (
	maxBodyLogSize = 64 * 1024 // 64KB

	// CtxAdminLogTitleKey 控制器可通过 c.Set(CtxAdminLogTitleKey, "标题") 自定义日志标题
	// 若未设置，中间件将自动从数据库查询规则标题
	CtxAdminLogTitleKey = "admin_log_title"
)

// AdminLog 管理员操作日志中间件
// 全局中间件，通过路径前缀判断是否记录
func AdminLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		fullPath := c.FullPath()
		if fullPath == "" {
			fullPath = c.Request.URL.Path
		}

		// 构建权限检查路径，用于匹配 AdminRule 的 name 字段
		adminPath := config.Get().Server.AdminBaseRoutePath
		checkPath := permission.BuildCheckPath(fullPath)

		// 仅记录管理后台请求
		// 仅记录 POST 请求
		// 额外排除 list 请求
		if !strings.HasPrefix(fullPath, adminPath+"/") || c.Request.Method != http.MethodPost || strings.HasSuffix(checkPath, "/list") {
			c.Next()
			return
		}

		// 捕获请求体，排除了 chunked 请求
		var bodyData string
		if c.Request.ContentLength > 0 {
			bodyBytes, err := io.ReadAll(io.LimitReader(c.Request.Body, maxBodyLogSize))

			if err != nil {
				httpx.Fail(c, httpx.WithMessage("操作日志记录失败: 读取 c.Request.Body 失败"))
				c.Abort()
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if len(bodyBytes) > 0 {
				bodyData = string(bodyBytes)
			}
		}

		c.Next()

		// 提前捕获所有需要的数据，避免 goroutine 中读取已复用的 context
		adminID := uint(0)
		adminName := ""
		admin := GetAdmin(c)
		if admin != nil {
			adminID = admin.ID
			adminName = admin.Username
		}

		clientIP := httpx.ClientIP(c)
		userAgent := util.TruncateStr(httpx.UserAgent(c), 255)

		// 检查控制器层是否已设置自定义日志标题
		customTitle, hasCustomTitle := c.Get(CtxAdminLogTitleKey)

		// 异步写入日志
		go func() {
			var title string
			repo := repoAuth.NewAdminLogRepository()
			if hasCustomTitle {
				title = customTitle.(string)
			} else {
				title = repo.GetRuleTitle(context.Background(), checkPath)
			}

			log := model.AdminLog{
				AdminID:   adminID,
				Username:  adminName,
				URL:       fullPath,
				Title:     title,
				IP:        clientIP,
				UserAgent: userAgent,
			}
			if bodyData != "" {
				log.Data = &bodyData
			}

			repo.Create(context.Background(), &log, repository.Options{})
		}()
	}
}
