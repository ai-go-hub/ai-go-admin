package handler

import (
	"strconv"

	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	"github.com/ai-go-hub/ai-go-admin/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler 通用控制器
type Handler[T any] struct {
	svc service.IService[T]
}

// NewHandler 创建通用控制器实例
func NewHandler[T any](svc service.IService[T]) *Handler[T] {
	return &Handler[T]{svc: svc}
}

// RegisterBaseRoutes 注册基控制器的通用 CRUD 路由
// 只使用 GET、POST，不使用 PUT、DELETE 等请求方式（除 GET/POST 外的请求方式，在国内的 CDN/全站加速 等场景兼容性极差）
func (h *Handler[T]) RegisterBaseRoutes(group *gin.RouterGroup) {
	group.POST("/create", h.Create)
	group.GET("/list", h.List)
	group.GET("/get/:id", h.Get)
	group.POST("/update/:id", h.Update)
	group.POST("/delete/:id", h.Delete)
}

// Create 新增记录
func (h *Handler[T]) Create(c *gin.Context) {
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	if err := h.svc.Create(c, &entity); err != nil {
		httpx.Fail(c, httpx.WithMessage("创建失败: "+err.Error()))
		return
	}

	httpx.Success(c)
}

// List 获取列表
func (h *Handler[T]) List(c *gin.Context) {
	entities, err := h.svc.List(c)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("查询失败: "+err.Error()))
		return
	}

	httpx.Success(c, httpx.WithData(entities))
}

// Get 获取单条记录
func (h *Handler[T]) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("无效的 ID"))
		return
	}

	entity, err := h.svc.Get(c, uint(id))
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("记录不存在"))
		return
	}

	httpx.Success(c, httpx.WithData(entity))
}

// Update 更新记录
func (h *Handler[T]) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("无效的 ID"))
		return
	}

	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	if err := h.svc.Update(c, uint(id), entity); err != nil {
		httpx.Fail(c, httpx.WithMessage("更新失败: "+err.Error()))
		return
	}

	httpx.Success(c)
}

// Delete 删除记录
func (h *Handler[T]) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("无效的 ID"))
		return
	}

	if err := h.svc.Delete(c, uint(id)); err != nil {
		httpx.Fail(c, httpx.WithMessage("删除失败: "+err.Error()))
		return
	}

	httpx.Success(c)
}
