package common

import (
	"github.com/ai-go-hub/ai-go-admin/internal/infra/upload"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"

	"github.com/gin-gonic/gin"
)

// UploadHandler 上传控制器
type UploadHandler struct {
	upload *upload.Service
}

// NewUploadHandler 创建上传控制器实例
func NewUploadHandler(up *upload.Service) *UploadHandler {
	return &UploadHandler{upload: up}
}

// Upload 处理文件上传
func (h *UploadHandler) Upload(c *gin.Context) {
	topic := c.PostForm("topic")
	driver := c.DefaultPostForm("driver", "local")

	file, err := c.FormFile("file")
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("没有文件被上传: "+err.Error()))
		return
	}

	// 获取当前登录用户信息，当前仅管理员一种，未登录禁止上传文件
	admin := middleware.GetAdmin(c)
	if admin == nil {
		httpx.Fail(c, httpx.WithMessage("请先登录"))
	}

	// 上传用户信息
	userId := admin.ID
	userType := "admin"

	res, err := h.upload.Upload(c.Request.Context(), file, driver, topic, userId, userType)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage(err.Error()))
		return
	}

	httpx.Success(c, httpx.WithData(res))
}

// RegisterRoutes 注册上传路由
func (h *UploadHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/upload", middleware.AdminAuthOptional(), h.Upload)
}
