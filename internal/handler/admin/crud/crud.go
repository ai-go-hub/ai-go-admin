package crud

import (
	"net/http"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/permission"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	svcCrud "github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"

	"github.com/gin-gonic/gin"
)

// Handler 可视化 CRUD 控制器
type Handler struct {
	svc *svcCrud.Service
}

// NewHandler 创建可视化 CRUD 控制器实例
func NewHandler(svc *svcCrud.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/model-list", middleware.AdminAuth(), h.ModelList)
	group.POST("/table-list", middleware.AdminAuth(), h.TableList)
	group.GET("/table-field-list", middleware.AdminAuth(), h.TableFieldList)
	group.GET("/generate-file-basic-data", middleware.AdminAuth(), h.GenerateFileBasicData)
}

// TableList 数据表列表
func (h *Handler) TableList(c *gin.Context) {
	if !h.checkPermission(c, []string{"crud/crud/read"}) {
		return
	}
	var req struct {
		ExcludeTables []string `json:"exclude_tables"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	tables, err := h.svc.TableList(c.Request.Context(), req.ExcludeTables)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("查询数据表失败: "+err.Error()))
		return
	}

	// 遍历修改 comment 为 "table - comment" 格式
	for i := range tables {
		tables[i].Comment = tables[i].Table + " - " + tables[i].Comment
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"list": tables,
	}))
}

// TableFieldList 数据表字段列表
func (h *Handler) TableFieldList(c *gin.Context) {
	if !h.checkPermission(c, []string{"crud/crud/read"}) {
		return
	}
	table := c.Query("table")
	if table == "" {
		httpx.Fail(c, httpx.WithMessage("参数错误: table 不能为空"))
		return
	}

	pk, fields, err := h.svc.TableFieldList(c.Request.Context(), table)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("查询数据表字段失败: "+err.Error()))
		return
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"pk":   pk,
		"list": fields,
	}))
}

// ModelList 模型列表
func (h *Handler) ModelList(c *gin.Context) {
	if !h.checkPermission(c, []string{"crud/crud/read"}) {
		return
	}
	models, err := h.svc.ModelList()
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("读取模型列表失败: "+err.Error()))
		return
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"list": models,
	}))
}

// GenerateFileBasicData 生成文件基本信息
func (h *Handler) GenerateFileBasicData(c *gin.Context) {
	if !h.checkPermission(c, []string{"crud/crud/create"}) {
		return
	}
	app := c.Query("app")
	table := c.Query("table")
	if table == "" || app == "" {
		httpx.Fail(c, httpx.WithMessage("参数错误: 不能为空"))
		return
	}

	types := []string{"model", "handler", "service", "repository", "router", "views"}
	data := make(map[string]svcCrud.GenerateFileBasicDataInfo, len(types))
	for _, typ := range types {
		data[typ] = svcCrud.GenerateFileBasicData(typ, table, app)
	}

	httpx.Success(c, httpx.WithData(data))
}

// checkPermission 校验 CRUD 权限节点，无权限时返回 false（内部已输出响应）
func (h *Handler) checkPermission(c *gin.Context, node []string) bool {
	admin := middleware.GetAdmin(c)
	if admin == nil {
		httpx.Fail(c, httpx.WithCode(http.StatusUnauthorized), httpx.WithMessage("未登录"))
		return false
	}

	perm := permission.New()
	hasPermission, err := perm.Check(c.Request.Context(), admin.ID, node, "AND")
	if err != nil || !hasPermission {
		httpx.Fail(c, httpx.WithMessage("无权限"))
		return false
	}
	return true
}
