package auth

import (
	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/handler"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"
	svcAuth "github.com/ai-go-hub/ai-go-admin/internal/service/admin/auth"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// AuthAdminHandler 管理员账号管理控制器
type AuthAdminHandler struct {
	*handler.Handler[model.Admin]
	svc *svcAuth.AuthAdminService
}

// NewAuthAdminHandler 创建管理员账号管理控制器实例
func NewAuthAdminHandler(svc *svcAuth.AuthAdminService) *AuthAdminHandler {
	return &AuthAdminHandler{
		Handler: handler.NewHandler(svc,
			handler.WithOmitFields(handler.ActionFields{
				// 创建时忽略以下字段不入库
				Create: []string{"id", "login_failure", "last_login_at", "last_login_ip", "deleted_at"},
			}),
			handler.WithPreloads([]repository.Preload{
				{Association: "AdminGroupAccesses.Group"},
			}),
			handler.WithExtension(func(c *gin.Context) any {
				return &svcAuth.AuthAdminExtension{
					// 避免 HTTP 层的中间件侵入到服务层，此处显式传递为扩展参数
					AdminSession: middleware.GetAdmin(c),
				}
			}),
		),
		svc: svc,
	}
}

// Create 覆写: 专用 DTO 接收请求数据（model.Admin.Password 有 json:"-"，无法直接绑定）
func (h *AuthAdminHandler) Create(c *gin.Context) {
	var req dto.AdminCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	entity := &model.Admin{}
	if err := copier.Copy(entity, req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数转换失败: "+err.Error()))
		return
	}

	if err := h.svc.Create(c, entity, h.BuildSerOpts(c, "Create", handler.Request{})); err != nil {
		httpx.Fail(c, httpx.WithMessage("创建失败: "+err.Error()))
		return
	}

	httpx.Success(c)
}

// Update 覆写: 用专用 DTO 接收请求数据，且 Password 为空时不更新
func (h *AuthAdminHandler) Update(c *gin.Context) {
	pk := c.Param("pk")

	var req dto.AdminUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	entity := model.Admin{}
	if err := copier.Copy(&entity, req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数转换失败: "+err.Error()))
		return
	}

	opts := h.BuildSerOpts(c, "Update", handler.Request{
		PrimaryKeyValue: pk,
	})

	if err := h.svc.Update(c, &entity, opts); err != nil {
		httpx.Fail(c, httpx.WithMessage("更新失败: "+err.Error()))
		return
	}

	httpx.Success(c)
}

// RegisterRoutes 注册路由
func (h *AuthAdminHandler) RegisterRoutes(group *gin.RouterGroup) {
	// 这种写法可自动挂载重写后的方法
	handler.RegisterBaseRoutes(h, group)
}
