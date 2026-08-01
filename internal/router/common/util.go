package common

import (
	handlerCommon "github.com/ai-go-hub/ai-go-admin/internal/handler/common"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.Register(func(r *gin.Engine) {
		h := handlerCommon.NewUtilHandler()

		group := r.Group("/common")
		h.RegisterRoutes(group)
	})
}
