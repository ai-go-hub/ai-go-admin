package routine

import (
	"strconv"

	handlerAuth "github.com/ai-go-hub/ai-go-admin/internal/handler/admin/auth"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
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

// Get 覆写: 校验 pk 是否为当前管理员 id
func (h *AuthAdminProfileHandler) Get(c *gin.Context) {
	pk := c.Param("pk")
	admin := middleware.GetAdmin(c)
	if admin == nil || strconv.FormatUint(uint64(admin.ID), 10) != pk {
		httpx.Fail(c, httpx.WithMessage("无权限"))
		return
	}
	h.AuthAdminHandler.Get(c)
}

// Update 覆写: 校验 pk 是否为当前管理员 id
func (h *AuthAdminProfileHandler) Update(c *gin.Context) {
	pk := c.Param("pk")
	admin := middleware.GetAdmin(c)
	if admin == nil || strconv.FormatUint(uint64(admin.ID), 10) != pk {
		httpx.Fail(c, httpx.WithMessage("无权限"))
		return
	}
	h.AuthAdminHandler.Update(c)
}
