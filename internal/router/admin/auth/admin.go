package auth

import (
	handlerAuth "github.com/ai-go-hub/ai-go-admin/internal/handler/admin/auth"
	repoAdmin "github.com/ai-go-hub/ai-go-admin/internal/repository/admin"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"
	svcAuth "github.com/ai-go-hub/ai-go-admin/internal/service/admin/auth"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.Register(func(r *gin.Engine) {
		repo := repoAdmin.NewAdminRepository()
		svc := svcAuth.NewAuthAdminService(repo)
		h := handlerAuth.NewAuthAdminHandler(svc)

		group := r.Group("/admin/auth/admin")
		h.RegisterRoutes(group)
	})
}
