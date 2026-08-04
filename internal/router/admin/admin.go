package admin

import (
	handlerAdmin "github.com/ai-go-hub/ai-go-admin/internal/handler/admin"
	repoAdmin "github.com/ai-go-hub/ai-go-admin/internal/repository/admin"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"
	svcAdmin "github.com/ai-go-hub/ai-go-admin/internal/service/admin"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.RegisterAdmin(func(group *gin.RouterGroup) {
		repo := repoAdmin.NewAdminRepository()
		svc := svcAdmin.NewAdminService(repo)
		h := handlerAdmin.NewAdminHandler(svc)

		h.RegisterRoutes(group)
	})
}
