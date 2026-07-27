package handler

import (
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	"github.com/ai-go-hub/ai-go-admin/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterBaseRoutes 注册通用 CRUD 路由
// 只使用 GET、POST，不使用 PUT、DELETE 等请求方式（除 GET/POST 外的请求方式，在国内的 CDN/全站加速 等场景兼容性极差）
func RegisterBaseRoutes(h IHandler, group *gin.RouterGroup) {
	group.POST("/list", middleware.AdminAuth(), h.List)
	group.POST("/create", middleware.AdminAuth(), h.Create)
	group.POST("/delete", middleware.AdminAuth(), h.Delete)
	group.POST("/sort", middleware.AdminAuth(), h.Sort)

	group.GET("/get/:pk", middleware.AdminAuth(), h.Get)
	group.POST("/update/:pk", middleware.AdminAuth(), h.Update)
}

// IHandler 通用控制器接口
type IHandler interface {
	Get(c *gin.Context)
	List(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Sort(c *gin.Context)
}

// Handler 通用控制器
type Handler[T any] struct {
	svc service.IService[T]
	cfg HandlerConfig
}

// Adapter 声明对应方法的数据适配器
type Adapter struct {
	Get  func(any, service.Options) (any, error)
	List func(any, service.Options) (any, error)
}

// ActionFields 声明对应方法的 Omit / Select 字段列表
type ActionFields struct {
	Get    []string
	List   []string
	Update []string
	Create []string
}

// ExtensionResolver 扩展参数解析函数
type ExtensionResolver func(c *gin.Context) any

// HandlerConfig 通用控制器配置
type HandlerConfig struct {
	Adapter           Adapter           // 数据适配器
	OmitFields        ActionFields      // 出入库黑名单字段
	SelectFields      ActionFields      // 出入库白名单字段
	ExtensionResolver ExtensionResolver // 任意自定义扩展数据的解析函数
}

// Option 通用控制器选项函数
type Option func(*HandlerConfig)

// Request 通用控制器请求体
type Request struct {
	Page             int                  `json:"page"`     // 页码
	Limit            int                  `json:"limit"`    // 每页条数
	Sort             string               `json:"sort"`     // 排序字段
	Order            string               `json:"order"`    // 排序方式
	Wheres           []service.WhereGroup `json:"wheres"`   // 查询条件组
	Selector         bool                 `json:"selector"` // 是否为来自选择器的查询
	PrimaryKeyValue  string               `json:"pk"`       // 主键值，前端传递简写 pk 即可
	PrimaryKeyValues []string             `json:"pks"`      // 主键值切片
}

// ListRequest 列表查询请求体
type ListRequest struct {
	Request
	PrimaryKeyValue  string   `json:"-"`
	PrimaryKeyValues []string `json:"-"`
}

// DeleteRequest 批量删除请求体
type DeleteRequest struct {
	PrimaryKeyValues []string `json:"pks" binding:"required,min=1"`
}

// SortRequest 拖动排序请求体
type SortRequest struct {
	Request
	Move      string `json:"move" binding:"required"`                    // 被拖动行主键
	Target    string `json:"target" binding:"required"`                  // 目标行主键
	Direction string `json:"direction" binding:"required,oneof=up down"` // 拖动方向
	Weigh     int64  `json:"weigh"`                                      // 目标行当前权重值
}

// WithAdapter 设置各动作的数据适配函数，可对数据进行额外加工
func WithAdapter(adapter Adapter) Option {
	return func(c *HandlerConfig) { c.Adapter = adapter }
}

// WithOmitFields 设置各动作的 Omit 黑名单字段
func WithOmitFields(f ActionFields) Option {
	return func(c *HandlerConfig) { c.OmitFields = f }
}

// WithSelectFields 设置各动作的 Select 白名单字段
func WithSelectFields(f ActionFields) Option {
	return func(c *HandlerConfig) { c.SelectFields = f }
}

// WithExtension 设置扩展参数解析函数
// 解析器返回值会被赋值给 service.Options.Extension
func WithExtension(resolve ExtensionResolver) Option {
	return func(c *HandlerConfig) { c.ExtensionResolver = resolve }
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

// Create 新增记录
func (h *Handler[T]) Create(c *gin.Context) {
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	opts := h.BuildSerOpts(c, "Create", Request{})
	if err := h.svc.Create(c, &entity, opts); err != nil {
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

	opts := h.BuildSerOpts(c, "Update", Request{
		PrimaryKeyValue: pk,
	})

	if err := h.svc.Update(c, &entity, opts); err != nil {
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

	opts := h.BuildSerOpts(c, "List", req.Request)
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

	resp := gin.H{
		"list":  list,
		"total": total,
	}

	if h.cfg.Adapter.List != nil {
		adapted, err := h.cfg.Adapter.List(list, opts)
		if err != nil {
			httpx.Fail(c, httpx.WithMessage("数据适配器错误: "+err.Error()))
			return
		}
		resp["list"] = adapted
	}

	httpx.Success(c, httpx.WithData(resp))
}

// Get 获取单条记录，按主键查询
func (h *Handler[T]) Get(c *gin.Context) {
	opts := h.BuildSerOpts(c, "Get", Request{
		PrimaryKeyValue: c.Param("pk"),
	})

	entity, err := h.svc.Get(c, opts)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("记录不存在"))
		return
	}

	resp := gin.H{
		"row": entity,
	}

	if h.cfg.Adapter.Get != nil {
		adapted, err := h.cfg.Adapter.Get(entity, opts)
		if err != nil {
			httpx.Fail(c, httpx.WithMessage("数据适配器错误: "+err.Error()))
			return
		}
		resp["row"] = adapted
	}

	httpx.Success(c, httpx.WithData(resp))
}

// Delete 批量删除记录
func (h *Handler[T]) Delete(c *gin.Context) {
	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	if err := h.svc.Delete(c, h.BuildSerOpts(c, "Delete", Request{
		PrimaryKeyValues: req.PrimaryKeyValues,
	})); err != nil {
		httpx.Fail(c, httpx.WithMessage("删除失败: "+err.Error()))
		return
	}

	httpx.Success(c)
}

// Sort 拖动排序
func (h *Handler[T]) Sort(c *gin.Context) {
	var req SortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	opts := h.BuildSerOpts(c, "Sort", req.Request)
	if err := h.svc.Sort(c, opts, req.Move, req.Target, req.Direction, req.Weigh); err != nil {
		httpx.Fail(c, httpx.WithMessage("调整排序失败: "+err.Error()))
		return
	}

	httpx.Success(c)
}

// BuildSerOpts 构建服务层选项数据
func (h *Handler[T]) BuildSerOpts(c *gin.Context, action string, opts Request) service.Options {
	serOpts := service.Options{
		Page:             opts.Page,
		Limit:            opts.Limit,
		SortField:        opts.Sort,
		SortOrder:        opts.Order,
		Wheres:           opts.Wheres,
		Selector:         opts.Selector,
		PrimaryKeyValue:  opts.PrimaryKeyValue,
		PrimaryKeyValues: opts.PrimaryKeyValues,
		OmitFields:       GetActionFields(action, h.cfg.OmitFields),
		SelectFields:     GetActionFields(action, h.cfg.SelectFields),
	}

	// 解析自定义扩展数据
	if h.cfg.ExtensionResolver != nil {
		serOpts.Extension = h.cfg.ExtensionResolver(c)
	}
	return serOpts
}

// GetActionFields 获取一个操作方法的 Omit / Select 字段列表
// 不使用反射，不使用 map[string][]string，此处简单粗暴的 switch 就是最好的
func GetActionFields(action string, fields ActionFields) []string {
	switch action {
	case "Get":
		return fields.Get
	case "List":
		return fields.List
	case "Update":
		return fields.Update
	case "Create":
		return fields.Create
	default:
		return nil
	}
}
