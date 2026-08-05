package routine

import (
	handlerAuth "github.com/ai-go-hub/ai-go-admin/internal/handler/admin/auth"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	svcAuth "github.com/ai-go-hub/ai-go-admin/internal/service/admin/auth"

	"github.com/gin-gonic/gin"
)

// AuthAdminProfileHandler 管理员个人资料控制器
// 嵌入 AuthAdminHandler，复用其 Get/Update 逻辑，仅覆写路由注册
type AuthAdminProfileHandler struct {
	*handlerAuth.AuthAdminHandler
}

// NewAuthAdminProfileHandler 创建管理员个人资料控制器
func NewAuthAdminProfileHandler(svc *svcAuth.AuthAdminService) *AuthAdminProfileHandler {
	return &AuthAdminProfileHandler{
		AuthAdminHandler: handlerAuth.NewAuthAdminHandler(svc),
	}
}

// RegisterRoutes 注册路由，只注册 get 和 update
func (h *AuthAdminProfileHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/get/:pk", middleware.AdminAuth(), middleware.AdminPermission(), h.Get)
	group.POST("/update/:pk", middleware.AdminAuth(), middleware.AdminPermission(), h.Update)
}
