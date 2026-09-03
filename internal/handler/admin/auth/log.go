package auth

import (
	"github.com/ai-go-hub/ai-go-admin/internal/handler"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	svcAdmin "github.com/ai-go-hub/ai-go-admin/internal/service/admin/auth"

	"github.com/gin-gonic/gin"
)

// AdminLogHandler 管理员日志控制器，嵌入通用控制器
type AdminLogHandler struct {
	*handler.Handler[model.AdminLog]
	svc *svcAdmin.AdminLogService
}

// NewAdminLogHandler 创建管理员日志控制器实例
func NewAdminLogHandler(svc *svcAdmin.AdminLogService) *AdminLogHandler {
	return &AdminLogHandler{
		Handler: handler.NewHandler(svc,
			handler.WithExtension(func(c *gin.Context) any {
				return &svcAdmin.AdminLogExtension{
					// 避免 HTTP 层的中间件侵入到服务层，此处显式传递为扩展参数
					AdminSession: middleware.GetAdmin(c),
				}
			}),
		),
		svc: svc,
	}
}

// RegisterRoutes 注册路由，日志只读，仅注册 List
// 无需验权，列表只会读取当前管理员的日志记录，同时菜单规则管理内保留下级 read 权限节点，用于决定菜单是否显示
func (h *AdminLogHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/list", middleware.AdminAuth(), h.List)
}
