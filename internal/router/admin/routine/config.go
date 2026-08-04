package routine

import (
	handlerRoutine "github.com/ai-go-hub/ai-go-admin/internal/handler/admin/routine"
	repoCommon "github.com/ai-go-hub/ai-go-admin/internal/repository/common"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"
	svcRoutine "github.com/ai-go-hub/ai-go-admin/internal/service/admin/routine"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.RegisterAdmin(func(group *gin.RouterGroup) {
		repo := repoCommon.NewConfigRepository()
		svc := svcRoutine.NewConfigService(repo)
		h := handlerRoutine.NewConfigHandler(svc, repo)

		subGroup := group.Group("/routine/config")
		h.RegisterRoutes(subGroup)
	})
}
