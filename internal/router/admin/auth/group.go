package auth

import (
	handlerAuth "github.com/ai-go-hub/ai-go-admin/internal/handler/admin/auth"
	repoAdmin "github.com/ai-go-hub/ai-go-admin/internal/repository/admin"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"
	svcAuth "github.com/ai-go-hub/ai-go-admin/internal/service/admin/auth"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.RegisterAdmin(func(group *gin.RouterGroup) {
		repo := repoAdmin.NewAdminGroupRepository()
		ruleRepo := repoAdmin.NewAdminRuleRepository()
		svc := svcAuth.NewAuthAdminGroupService(repo, ruleRepo)
		h := handlerAuth.NewAuthAdminGroupHandler(svc)

		subGroup := group.Group("/auth/group")
		h.RegisterRoutes(subGroup)
	})
}
