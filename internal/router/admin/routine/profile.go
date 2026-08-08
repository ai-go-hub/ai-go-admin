package routine

import (
	handlerAuth "github.com/ai-go-hub/ai-go-admin/internal/handler/admin/routine"
	repoAdmin "github.com/ai-go-hub/ai-go-admin/internal/repository/admin"
	repoAuth "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/auth"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"
	svcAuth "github.com/ai-go-hub/ai-go-admin/internal/service/admin/auth"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.RegisterAdmin(func(group *gin.RouterGroup) {
		repo := repoAdmin.NewAdminRepository()
		groupRepo := repoAuth.NewAdminGroupRepository()
		svc := svcAuth.NewAuthAdminService(repo, groupRepo)
		h := handlerAuth.NewAuthAdminProfileHandler(svc)

		subGroup := group.Group("/routine/profile")
		h.RegisterRoutes(subGroup)
	})
}
