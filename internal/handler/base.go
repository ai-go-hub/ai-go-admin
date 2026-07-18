package handler

import (
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	"github.com/ai-go-hub/ai-go-admin/internal/service"

	"github.com/gin-gonic/gin"
)

// IHandler 通用控制器接口
type IHandler interface {
	Get(c *gin.Context)
	List(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// Handler 通用控制器
type Handler[T any] struct {
	svc service.IService[T]
	cfg HandlerConfig
}

// ListRequest 列表查询请求体
type ListRequest struct {
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
	Sort   string          `json:"sort"`
	Order  string          `json:"order"`
	Wheres []service.Where `json:"wheres"`
}

// OmitFields 声明各方法的 Omit 黑名单字段，留空代表不忽略任何字段
type OmitFields struct {
	Get    []string
	List   []string
	Update []string
	Create []string
}

// SelectFields 声明各方法的 Select 白名单字段，留空代表全部字段
type SelectFields struct {
	Get    []string
	List   []string
	Update []string
	Create []string
}

// HandlerConfig 通用控制器配置
type HandlerConfig struct {
	Omit   OmitFields
	Select SelectFields
}

// Option 通用控制器选项函数
type Option func(*HandlerConfig)

// WithOmitFields 设置各动作的 Omit 黑名单字段
func WithOmitFields(f OmitFields) Option {
	return func(c *HandlerConfig) { c.Omit = f }
}

// WithSelectFields 设置各动作的 Select 白名单字段
func WithSelectFields(f SelectFields) Option {
	return func(c *HandlerConfig) { c.Select = f }
}

// NewHandler 创建通用控制器实例，支持函数式选项传递配置
func NewHandler[T any](svc service.IService[T], opts ...Option) *Handler[T] {
	cfg := HandlerConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Handler[T]{svc: svc, cfg: cfg}
}

// Config 返回控制器配置，供嵌入子控制器读取
func (h *Handler[T]) Config() HandlerConfig {
	return h.cfg
}

// RegisterBaseRoutes 注册通用 CRUD 路由
// 只使用 GET、POST，不使用 PUT、DELETE 等请求方式（除 GET/POST 外的请求方式，在国内的 CDN/全站加速 等场景兼容性极差）
func RegisterBaseRoutes(h IHandler, group *gin.RouterGroup) {
	group.POST("/list", h.List, middleware.AdminAuth())
	group.POST("/create", h.Create, middleware.AdminAuth())

	group.GET("/get/:pk", h.Get, middleware.AdminAuth())
	group.POST("/update/:pk", h.Update, middleware.AdminAuth())
	group.POST("/delete/:pk", h.Delete, middleware.AdminAuth())
}

// Create 新增记录
func (h *Handler[T]) Create(c *gin.Context) {
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	if err := h.svc.Create(c, &entity, service.Options{
		Omit:   h.cfg.Omit.Create,
		Select: h.cfg.Select.Create,
	}); err != nil {
		httpx.Fail(c, httpx.WithMessage("创建失败: "+err.Error()))
		return
	}

	httpx.Success(c)
}

// Update 更新记录
func (h *Handler[T]) Update(c *gin.Context) {
	pk := c.Param("pk")

	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	if err := h.svc.Update(c, pk, entity, service.Options{
		Omit:   h.cfg.Omit.Update,
		Select: h.cfg.Select.Update,
	}); err != nil {
		httpx.Fail(c, httpx.WithMessage("更新失败: "+err.Error()))
		return
	}

	httpx.Success(c)
}

// List 获取数据列表
func (h *Handler[T]) List(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	opts := service.Options{
		Omit:      h.cfg.Omit.List,
		Select:    h.cfg.Select.List,
		Wheres:    req.Wheres,
		SortField: req.Sort,
		SortOrder: req.Order,
		Page:      req.Page,
		Limit:     req.Limit,
	}
	list, err := h.svc.List(c, opts)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("查询失败: "+err.Error()))
		return
	}

	total, err := h.svc.Count(c, opts)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("总数查询失败: "+err.Error()))
		return
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"list":  list,
		"total": total,
	}))
}

// Get 获取单条记录，按主键查询
func (h *Handler[T]) Get(c *gin.Context) {
	entity, err := h.svc.Get(c, service.Options{
		Omit:       h.cfg.Omit.Get,
		Select:     h.cfg.Select.Get,
		PrimaryKey: c.Param("pk"),
	})
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("记录不存在"))
		return
	}

	httpx.Success(c, httpx.WithData(entity))
}

// Delete 删除记录
func (h *Handler[T]) Delete(c *gin.Context) {
	if err := h.svc.Delete(c, c.Param("pk")); err != nil {
		httpx.Fail(c, httpx.WithMessage("删除失败: "+err.Error()))
		return
	}

	httpx.Success(c)
}
