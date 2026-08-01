package routine

import (
	handlerRoutine "github.com/ai-go-hub/ai-go-admin/internal/handler/admin/routine"
	repoCommon "github.com/ai-go-hub/ai-go-admin/internal/repository/common"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"
	svcRoutine "github.com/ai-go-hub/ai-go-admin/internal/service/admin/routine"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.Register(func(r *gin.Engine) {
		repo := repoCommon.NewAttachmentRepository()
		svc := svcRoutine.NewAttachmentService(repo)
		h := handlerRoutine.NewAttachmentHandler(svc)

		group := r.Group("/admin/routine/attachment")
		h.RegisterRoutes(group)
	})
}
