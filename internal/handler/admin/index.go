package admin

import (
	"ai-go-mall/internal/handler"
	"ai-go-mall/internal/model"
	"ai-go-mall/internal/response"
	svcAdmin "ai-go-mall/internal/service/admin"

	"github.com/gin-gonic/gin"
)

// Handler 管理员控制器，嵌入通用控制器并扩展自定义方法
type Handler struct {
	*handler.Handler[model.Admin]
	svc *svcAdmin.Service
}

// NewHandler 创建管理员控制器实例
func NewHandler(svc *svcAdmin.Service) *Handler {
	return &Handler{
		Handler: handler.NewHandler(svc),
		svc:     svc,
	}
}

// Login 管理员登录
func (h *Handler) Login(c *gin.Context) {
	var req svcAdmin.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, response.WithMessage("参数错误: "+err.Error()))
		return
	}

	result, err := h.svc.Login(c, &req)
	if err != nil {
		response.Fail(c, response.WithMessage(err.Error()))
		return
	}

	response.Success(c, response.WithData(result))
}

// RegisterRoutes 注册路由（包含基控制器路由 + 自定义路由）
func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	// 只注册自定义路由
	// 不注册基控制器的 CRUD 路由
	group.POST("/login", h.Login)
}
