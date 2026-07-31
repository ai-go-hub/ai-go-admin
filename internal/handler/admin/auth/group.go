package auth

import (
	"github.com/ai-go-hub/ai-go-admin/internal/handler"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	svcAuth "github.com/ai-go-hub/ai-go-admin/internal/service/admin/auth"
	"github.com/ai-go-hub/ai-go-admin/pkg/tree"

	"github.com/gin-gonic/gin"
)

// AuthAdminGroupHandler 管理员分组管理控制器
type AuthAdminGroupHandler struct {
	*handler.Handler[model.AdminGroup]
	svc *svcAuth.AuthAdminGroupService
}

// NewAuthAdminGroupHandler 创建管理员分组管理控制器实例
func NewAuthAdminGroupHandler(svc *svcAuth.AuthAdminGroupService) *AuthAdminGroupHandler {
	return &AuthAdminGroupHandler{
		Handler: handler.NewHandler(svc,
			handler.WithExtension(func(c *gin.Context) any {
				return &svcAuth.AuthAdminGroupExtension{
					// 避免 HTTP 层的中间件侵入到服务层，此处显式传递为扩展参数
					AdminSession: middleware.GetAdmin(c),
				}
			}),
		),
		svc: svc,
	}
}

// List 覆写通用列表方法: 追加 rules_title 摘要字段，并渲染为树状结构
func (h *AuthAdminGroupHandler) List(c *gin.Context) {
	var req handler.ListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	opts := h.BuildSerOpts(c, "List", req.Request)
	groups, err := h.svc.List(c, opts)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("查询失败: "+err.Error()))
		return
	}

	total, err := h.svc.Count(c, opts)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("总数查询失败: "+err.Error()))
		return
	}

	// 批量构建每个分组的 rules 摘要文本（例如: 控制台等 60 项）
	titles, err := h.svc.BuildRulesTitles(c, groups)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("规则标题查询失败: "+err.Error()))
		return
	}

	// 转 map 以备后续树状数据组装，并注入 rules_title
	groupData := make([]map[string]any, len(groups))
	for i := range groups {
		m := groups[i].ToMap()
		m["rules_title"] = titles[groups[i].ID]
		groupData[i] = m
	}

	var list []map[string]any
	if opts.Selector {
		list = tree.Render(groupData, "id", "pid", "name")
	} else {
		list = tree.Build(groupData, "id", "pid", "children")
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"list":  list,
		"total": total,
	}))
}

// Get 覆写通用单条读取方法: 出库后剥离 rules 中的父级 ID（父级选择状态由子级决定）
func (h *AuthAdminGroupHandler) Get(c *gin.Context) {
	opts := h.BuildSerOpts(c, "Get", handler.Request{
		PrimaryKeyValue: c.Param("pk"),
	})

	entity, err := h.svc.Get(c, opts)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("记录不存在"))
		return
	}

	stripped, err := h.svc.StripParentRuleIDs(c, entity)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("规则数据处理失败: "+err.Error()))
		return
	}
	if stripped != nil {
		entity.Rules = stripped
	}

	httpx.Success(c, httpx.WithData(gin.H{"row": entity}))
}

// RegisterRoutes 注册路由
func (h *AuthAdminGroupHandler) RegisterRoutes(group *gin.RouterGroup) {
	handler.RegisterBaseRoutes(h, group)
}
