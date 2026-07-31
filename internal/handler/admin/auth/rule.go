package auth

import (
	"github.com/ai-go-hub/ai-go-admin/internal/handler"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/service"
	svcAuth "github.com/ai-go-hub/ai-go-admin/internal/service/admin/auth"
	"github.com/ai-go-hub/ai-go-admin/pkg/tree"

	"github.com/gin-gonic/gin"
)

// AuthAdminRuleHandler 菜单和权限规则管理控制器
type AuthAdminRuleHandler struct {
	*handler.Handler[model.AdminRule]
	svc *svcAuth.AuthAdminRuleService
}

// NewAuthAdminRuleHandler 创建菜单和权限规则管理控制器实例
func NewAuthAdminRuleHandler(svc *svcAuth.AuthAdminRuleService) *AuthAdminRuleHandler {
	return &AuthAdminRuleHandler{
		Handler: handler.NewHandler(svc,
			handler.WithAdapter(handler.Adapter{
				// 定义控制器层的数据适配器（就不需要重写控制器层的 List 方法了）
				List: func(data any, opts service.Options) (any, error) {
					rules, ok := data.([]model.AdminRule)
					if !ok {
						return data, nil
					}

					ruleData := make([]map[string]any, len(rules))
					for i := range rules {
						ruleData[i] = rules[i].ToMap()
					}

					if opts.Selector {
						return tree.Render(ruleData, "id", "pid", "title"), nil
					} else {
						return tree.Build(ruleData, "id", "pid", "children"), nil
					}
				},
			}),
			handler.WithExtension(func(c *gin.Context) any {
				return &svcAuth.AuthAdminRuleExtension{
					// 避免 HTTP 层的中间件侵入到服务层，此处显式传递为扩展参数
					AdminSession: middleware.GetAdmin(c),
				}
			}),
		),
		svc: svc,
	}
}

// RegisterRoutes 注册路由
func (h *AuthAdminRuleHandler) RegisterRoutes(group *gin.RouterGroup) {
	handler.RegisterBaseRoutes(h, group)
}
