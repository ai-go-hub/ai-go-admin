package crud

import (
	"errors"
	"net/http"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/permission"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	repoAuth "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/auth"
	repoCrud "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/crud"
	svcCrud "github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"
	"github.com/ai-go-hub/ai-go-admin/pkg/filesystem"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 可视化 CRUD 控制器
type Handler struct {
	svc      *svcCrud.Service
	repo     *repoCrud.CrudLogRepository
	repoAuth *repoAuth.AdminRuleRepository
}

// NewHandler 创建可视化 CRUD 控制器实例
func NewHandler(svc *svcCrud.Service, repo *repoCrud.CrudLogRepository, repoAuth *repoAuth.AdminRuleRepository) *Handler {
	return &Handler{
		svc:      svc,
		repo:     repo,
		repoAuth: repoAuth,
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/table-field-list", middleware.AdminAuth(), h.TableFieldList)
	group.GET("/generate-basic-data", middleware.AdminAuth(), h.GenerateBasicData)
	group.GET("/check-log", middleware.AdminAuth(), h.CheckLog)
	group.GET("/parse-table-data", middleware.AdminAuth(), h.ParseTableData)

	group.POST("/model-list", middleware.AdminAuth(), h.ModelList)
	group.POST("/table-list", middleware.AdminAuth(), h.TableList)
	group.POST("/check-generate", middleware.AdminAuth(), h.CheckGenerate)
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

	// 遍历修改 comment 为 "model name - comment" 格式
	for i := range models {
		models[i].Comment = models[i].Name + " - " + models[i].Comment
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"list": models,
	}))
}

// GenerateBasicData 生成 CRUD 基本信息（文件生成位置、路由路径）
func (h *Handler) GenerateBasicData(c *gin.Context) {
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
	files := make(map[string]svcCrud.GenerateFileBasicDataInfo, len(types))
	for _, typ := range types {
		files[typ] = svcCrud.GenerateFileBasicData(typ, table, app)
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"files": files,
		"route": svcCrud.GenerateRoutePath(table),
	}))
}

// CheckLog 查询指定数据表的 CRUD 记录
func (h *Handler) CheckLog(c *gin.Context) {
	if !h.checkPermission(c, []string{"crud/crud/read"}) {
		return
	}
	table := c.Query("table")
	if table == "" {
		httpx.Fail(c, httpx.WithMessage("参数错误: table 不能为空"))
		return
	}

	log, err := h.repo.FindSucceededByName(c.Request.Context(), table)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("查询 CRUD 记录失败: "+err.Error()))
		return
	}

	var id any
	if log != nil {
		id = log.ID
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"id": id,
	}))
}

// CheckGenerate 生成前检查
func (h *Handler) CheckGenerate(c *gin.Context) {
	var req struct {
		Table       string `json:"table"`
		HandlerFile string `json:"handler"`
		RoutePath   string `json:"route"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	// 检查 handler 文件是否存在（存在则生成时会覆盖）
	fileExists, err := filesystem.Exists(req.HandlerFile)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("检查控制器文件失败: "+err.Error()))
		return
	}

	// 检查数据表是否存在
	tableExists, err := h.svc.TableExists(c.Request.Context(), req.Table)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("检查数据表失败: "+err.Error()))
		return
	}

	// 检查菜单规则是否存在（未查得时 menu 返回 false，不报错）
	menu, err := h.repoAuth.FindByPath(c.Request.Context(), req.RoutePath)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		httpx.Fail(c, httpx.WithMessage("检查菜单规则失败: "+err.Error()))
		return
	}
	menuExists := menu != nil

	if !fileExists || !tableExists || !menuExists {
		httpx.Fail(c, httpx.WithData(gin.H{
			"handler": fileExists,
			"table":   tableExists,
			"menu":    menuExists,
		}), httpx.WithCode(-1))
		return
	}

	httpx.Success(c)
}

// ParseTableData 解析指定数据表的字段数据
func (h *Handler) ParseTableData(c *gin.Context) {
	if !h.checkPermission(c, []string{"crud/crud/read"}) {
		return
	}

	table := c.Query("table")
	if table == "" {
		httpx.Fail(c, httpx.WithMessage("参数错误: table 不能为空"))
		return
	}

	columns, err := h.svc.ParseFieldData(c.Request.Context(), table)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("解析字段数据失败: "+err.Error()))
		return
	}

	comment, err := svcCrud.TableComment(c.Request.Context(), table)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("查询数据表注释失败: "+err.Error()))
		return
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"columns": columns,
		"comment": comment,
	}))
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
