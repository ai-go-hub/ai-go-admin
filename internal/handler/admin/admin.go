package admin

import (
	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/handler"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/response"
	svcAdmin "github.com/ai-go-hub/ai-go-admin/internal/service/admin"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理员控制器，嵌入通用控制器并扩展自定义方法
type AdminHandler struct {
	*handler.Handler[model.Admin]
	svc *svcAdmin.AdminService
}

// NewAdminHandler 创建管理员控制器实例
func NewAdminHandler(svc *svcAdmin.AdminService) *AdminHandler {
	return &AdminHandler{
		Handler: handler.NewHandler(svc),
		svc:     svc,
	}
}

// GetLoginConfig 返回管理员登录页配置（供前端判断是否启用人机验证码）
func (h *AdminHandler) GetLoginConfig(c *gin.Context) {
	admin := middleware.GetAdmin(c)
	flagType := middleware.FlagNeedLogin
	if admin != nil {
		flagType = middleware.FlagLoggedIn
	}

	response.Success(c, response.WithData(gin.H{
		"type":    flagType,
		"captcha": config.Get().Captcha.Switches.AdminLogin,
	}))
}

// Login 管理员登录
func (h *AdminHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
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

// Logout 管理员退出登录
func (h *AdminHandler) Logout(c *gin.Context) {
	admin := middleware.GetAdmin(c)
	if admin == nil {
		// 未登录，直接生成成功响应，其意图已自然完成
		// 不能执行任何 token 删除逻辑，因为管理员未认证
		response.Success(c)
		return
	}

	if err := h.svc.Logout(c, middleware.GetToken(c)); err != nil {
		response.Fail(c)
		return
	}

	response.Success(c)
}

// ClearCache 清理后台缓存
//
// 当前后端无可清理的缓存（暂无 Redis/进程内运行时缓存层等）
// 此接口作为契约占位，待后续引入缓存层时在此补充清理逻辑
func (h *AdminHandler) ClearCache(c *gin.Context) {
	response.Success(c)
}

// Init 后台数据初始化接口
func (h *AdminHandler) Init(c *gin.Context) {
	result, err := h.svc.Init(c, middleware.GetAdmin(c))
	if err != nil {
		response.Fail(c, response.WithMessage("初始化失败: "+err.Error()))
		return
	}

	response.Success(c, response.WithData(result))
}

// RegisterRoutes 注册路由
func (h *AdminHandler) RegisterRoutes(group *gin.RouterGroup) {
	// 只注册自定义路由
	// 不注册基控制器的 CRUD 路由
	group.POST("/login", h.Login)
	group.POST("/logout", middleware.AdminAuthOptional(), h.Logout)
	group.GET("/login-config", middleware.AdminAuthOptional(), h.GetLoginConfig)

	group.GET("/init", middleware.AdminAuth(), h.Init)
	group.POST("/clear-cache", middleware.AdminAuth(), h.ClearCache)
}
