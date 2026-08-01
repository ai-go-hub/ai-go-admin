package auth

import (
	handlerAdmin "github.com/ai-go-hub/ai-go-admin/internal/handler/admin"
	repoAdmin "github.com/ai-go-hub/ai-go-admin/internal/repository/admin"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"
	svcAdmin "github.com/ai-go-hub/ai-go-admin/internal/service/admin"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.Register(func(r *gin.Engine) {
		repo := repoAdmin.NewAdminLogRepository()
		svc := svcAdmin.NewAdminLogService(repo)
		h := handlerAdmin.NewAdminLogHandler(svc)

		group := r.Group("/admin/auth/admin-log")
		h.RegisterRoutes(group)
	})
}
