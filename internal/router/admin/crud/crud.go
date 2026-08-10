package crud

import (
	handlerCrud "github.com/ai-go-hub/ai-go-admin/internal/handler/admin/crud"
	repoAuth "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/auth"
	repoCrud "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/crud"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"
	svcCrud "github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.RegisterAdmin(func(group *gin.RouterGroup) {
		svc := svcCrud.NewService()
		h := handlerCrud.NewHandler(svc, repoCrud.NewCrudLogRepository(), repoAuth.NewAdminRuleRepository())

		subGroup := group.Group("/crud")
		h.RegisterRoutes(subGroup)
	})
}
