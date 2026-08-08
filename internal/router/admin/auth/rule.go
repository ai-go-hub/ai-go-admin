package auth

import (
	handlerAuth "github.com/ai-go-hub/ai-go-admin/internal/handler/admin/auth"
	repoAuth "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/auth"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"
	svcAuth "github.com/ai-go-hub/ai-go-admin/internal/service/admin/auth"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.RegisterAdmin(func(group *gin.RouterGroup) {
		repo := repoAuth.NewAdminRuleRepository()
		svc := svcAuth.NewAuthAdminRuleService(repo)
		h := handlerAuth.NewAuthAdminRuleHandler(svc)

		subGroup := group.Group("/auth/rule")
		h.RegisterRoutes(subGroup)
	})
}
