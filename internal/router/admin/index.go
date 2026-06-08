package admin

import (
	handlerAdmin "ai-go-mall/internal/handler/admin"
	repoAdmin "ai-go-mall/internal/repository/admin"
	"ai-go-mall/internal/router/registry"
	svcAdmin "ai-go-mall/internal/service/admin"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.Register(func(r *gin.Engine) {
		repo := repoAdmin.NewRepository()
		svc := svcAdmin.NewService(repo)
		h := handlerAdmin.NewHandler(svc)

		group := r.Group("/admin")
		h.RegisterRoutes(group)
	})
}
