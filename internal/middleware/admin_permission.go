package middleware

import (
	"net/http"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/permission"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"

	"github.com/gin-gonic/gin"
)

// AdminPermission 管理员权限校验中间件
// 必须在 AdminAuth（管理员认证中间件）之后执行，仅挂载到需要验权的路由上，无需验权的路由不要注册此中间件
func AdminPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		admin := GetAdmin(c)
		if admin == nil {
			httpx.Fail(c, httpx.WithCode(http.StatusUnauthorized), httpx.WithMessage("未登录"))
			c.Abort()
			return
		}

		// 构建检查路径: 去路径中的 query 和 params，去前缀
		fullPath := c.FullPath()
		if fullPath == "" {
			fullPath = c.Request.URL.Path
		}
		checkPath := permission.BuildCheckPath(fullPath)

		// 构建验权参数
		op := "AND"
		ruleNames := []string{checkPath}

		// list 和 get 同时验 read，op 为 OR
		if lastSlash := strings.LastIndex(checkPath, "/"); lastSlash != -1 {
			lastSeg := checkPath[lastSlash+1:]
			if lastSeg == "list" || lastSeg == "get" {
				ruleNames = []string{checkPath, checkPath[:lastSlash] + "/read"}
				op = "OR"
			}
		}

		perm := permission.New()
		hasPermission, err := perm.Check(c.Request.Context(), admin.ID, ruleNames, op)
		if err != nil || !hasPermission {
			httpx.Fail(c, httpx.WithCode(http.StatusForbidden), httpx.WithMessage("无权限"))
			c.Abort()
			return
		}

		c.Next()
	}
}
