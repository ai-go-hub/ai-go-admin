package routine

import (
	"github.com/ai-go-hub/ai-go-admin/internal/handler"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	svcRoutine "github.com/ai-go-hub/ai-go-admin/internal/service/admin/routine"

	"github.com/gin-gonic/gin"
)

// AttachmentHandler 附件控制器
type AttachmentHandler struct {
	*handler.Handler[model.Attachment]
	svc *svcRoutine.AttachmentService
}

// NewAttachmentHandler 创建附件控制器实例
func NewAttachmentHandler(svc *svcRoutine.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{
		Handler: handler.NewHandler(svc,
			handler.WithOmitFields(handler.ActionFields{
				Create: []string{"id", "url", "driver"},
				Update: []string{"id", "url", "driver"},
			}),
		),
		svc: svc,
	}
}

// RegisterRoutes 注册路由
func (h *AttachmentHandler) RegisterRoutes(group *gin.RouterGroup) {
	handler.RegisterBaseRoutes(h, group)
}
