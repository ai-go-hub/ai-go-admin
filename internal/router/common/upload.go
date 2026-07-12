package common

import (
	handlerCommon "github.com/ai-go-hub/ai-go-admin/internal/handler/common"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/upload"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.Register(func(r *gin.Engine) {
		h := handlerCommon.NewUploadHandler(upload.Manager())

		group := r.Group("/common")
		h.RegisterRoutes(group)
	})
}
