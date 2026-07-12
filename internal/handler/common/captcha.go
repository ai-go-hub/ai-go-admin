package common

import (
	"github.com/ai-go-hub/ai-go-admin/internal/infra/captcha"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	svcCommon "github.com/ai-go-hub/ai-go-admin/internal/service/common"

	"github.com/gin-gonic/gin"
)

// CaptchaHandler 验证码控制器
type CaptchaHandler struct {
	svc *svcCommon.CaptchaService
}

// NewCaptchaHandler 创建验证码控制器实例
func NewCaptchaHandler(svc *svcCommon.CaptchaService) *CaptchaHandler {
	return &CaptchaHandler{svc: svc}
}

// Create 创建点选验证码
func (h *CaptchaHandler) Create(c *gin.Context) {
	result, err := h.svc.Create()
	if err != nil {
		httpx.Fail(c, httpx.WithMessage(err.Error()))
		return
	}
	httpx.Success(c, httpx.WithData(result))
}

// Verify 校验点选验证码，仅用于预验场景，校验成功不会使验证码失效
func (h *CaptchaHandler) Verify(c *gin.Context) {
	var req captcha.Request
	if err := c.ShouldBind(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	_, err := h.svc.Verify(req)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage(err.Error()))
		return
	}
	httpx.Success(c)
}

// RegisterRoutes 注册验证码路由
func (h *CaptchaHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/create", h.Create)
	group.POST("/verify", h.Verify)
}
