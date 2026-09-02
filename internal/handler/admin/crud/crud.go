package crud

import (
	"bufio"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/permission"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"
	repoAuth "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/auth"
	repoCrud "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/crud"
	svcCrud "github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"
	svcGen "github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud/gen"
	tbl "github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud/table"
	"github.com/ai-go-hub/ai-go-admin/pkg/airx"
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
	group.GET("/log-start", middleware.AdminAuth(), h.LogStart)

	group.POST("/model-list", middleware.AdminAuth(), h.ModelList)
	group.POST("/table-list", middleware.AdminAuth(), h.TableList)
	group.POST("/check-generate", middleware.AdminAuth(), h.CheckGenerate)
	group.POST("/generate", middleware.AdminAuth(), h.Generate)

	group.POST("/ai/stream", middleware.AdminAuth(), h.ChatStream)
}

// TableList 数据表列表
func (h *Handler) TableList(c *gin.Context) {
	if !h.checkPermission(c, []string{"crud/crud/read"}) {
		return
	}
	var req struct {
		Exclusions []string `json:"exclusions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	tables, err := h.svc.TableList(c.Request.Context(), req.Exclusions)
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
	var req struct {
		Exclusions []string `json:"exclusions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	models, err := h.svc.ModelList(req.Exclusions)
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
	path := c.Query("path")
	if path == "" || app == "" {
		httpx.Fail(c, httpx.WithMessage("参数错误: 不能为空"))
		return
	}

	types := []string{"model", "handler", "service", "repository", "router", "views"}
	files := make(map[string]dto.GenerateFileBasicDataInfo, len(types))
	for _, typ := range types {
		files[typ] = svcCrud.GenerateFileBasicData(typ, path, app)
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"files": files,
		"route": svcCrud.GenerateRoutePath(path),
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
	if !h.checkPermission(c, []string{"crud/crud/create"}) {
		return
	}
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

	// 判断对应数据表是否为空
	empty, err := h.svc.TableEmpty(c.Request.Context(), table)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("查询数据表是否为空失败: "+err.Error()))
		return
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"columns": columns,
		"comment": comment,
		"empty":   empty,
	}))
}

// LogStart 从生成记录开始
func (h *Handler) LogStart(c *gin.Context) {
	if !h.checkPermission(c, []string{"crud/crud/read"}) {
		return
	}
	id := c.Query("id")
	if id == "" {
		httpx.Fail(c, httpx.WithMessage("参数错误: 记录 id 不能为空"))
		return
	}

	log, err := h.repo.Get(c.Request.Context(), repository.Options{
		PrimaryKeyValue: id,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpx.Fail(c, httpx.WithMessage("CRUD 记录不存在"))
			return
		}
		httpx.Fail(c, httpx.WithMessage("查询 CRUD 记录失败: "+err.Error()))
		return
	}

	// 判断对应数据表是否为空
	empty, err := h.svc.TableEmpty(c.Request.Context(), log.Name)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("查询数据表是否为空失败: "+err.Error()))
		return
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"fields": log.Fields,
		"table":  log.Table,
		"empty":  empty,
	}))
}

// Generate 生成 CRUD 代码
func (h *Handler) Generate(c *gin.Context) {
	if !h.checkPermission(c, []string{"crud/crud/create"}) {
		return
	}
	var req struct {
		Type   string           `json:"type"`
		Table  dto.CRUDTable    `json:"table"`
		Fields []dto.CRUDFields `json:"fields"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	// 生成期间暂停 air 热更新，避免写入代码触发重编译重启导致本请求被中断
	airx.Pause()
	defer airx.Resume()

	// 记录生成日志
	logID, err := svcCrud.CrudLog(c.Request.Context(), h.repo, dto.CrudLogData{
		Table:  req.Table,
		Fields: req.Fields,
		Status: "generating",
	})
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("记录生成日志失败: "+err.Error()))
		return
	}

	// CRUD生成失败时更新日志状态为 failed，并输出失败响应
	failLog := func(msg string) {
		if logID > 0 {
			_ = svcCrud.UpdateLogStatus(c.Request.Context(), h.repo, logID, "failed")
		}
		httpx.Fail(c, httpx.WithMessage(msg))
	}

	// 待生成文件列表
	generateFileMap := map[string]string{
		"model":      req.Table.ModelFile,
		"handler":    req.Table.HandlerFile,
		"service":    req.Table.ServiceFile,
		"repository": req.Table.RepositoryFile,
		"router":     req.Table.RouterFile,
		"views":      req.Table.WebViewsDir,
	}

	// 根据前端传入的文件路径，解析和组装各生成文件的基础数据
	generateFileBasicData := make(map[string]dto.GenerateFileBasicDataInfo, len(generateFileMap)+1)
	for name, path := range generateFileMap {
		generateFileBasicData[name] = svcCrud.ParseGenerateFileBasicData(name, path)
	}
	generateFileBasicData["lang"] = svcCrud.GenerateFileBasicData("lang", generateFileBasicData["views"].Path, generateFileBasicData["views"].App)

	// 按需删除原表
	if req.Type == "create" || req.Table.Rebuild == "Yes" {
		if err := tbl.DropTable(c.Request.Context(), svcCrud.WithPrefix(req.Table.Name)); err != nil {
			failLog("删除数据表失败: " + err.Error())
			return
		}
	}

	// 同步数据表结构变更，并收集执行的 SQL
	sqls, err := svcCrud.HandleTableDesign(c.Request.Context(), req.Table, req.Fields)
	if err != nil {
		failLog("同步数据表结构失败: " + err.Error())
		return
	}
	// 记录执行的 SQL 到 CRUD 日志
	if logID > 0 && len(sqls) > 0 {
		sqlText := strings.Join(sqls, "\n")
		if _, err := svcCrud.CrudLog(c.Request.Context(), h.repo, dto.CrudLogData{
			ID:  &logID,
			Sql: &sqlText,
		}); err != nil {
			failLog("记录执行的 SQL 失败: " + err.Error())
			return
		}
	}

	// 创建模型文件
	if err := svcGen.CreateModelFile(generateFileBasicData, req.Table, req.Fields); err != nil {
		failLog("生成模型文件失败: " + err.Error())
		return
	}

	// 创建仓储、服务、控制器、路由文件
	if err := svcGen.CreateRepositoryFile(generateFileBasicData, req.Table, req.Fields); err != nil {
		failLog("生成仓储文件失败: " + err.Error())
		return
	}
	if err := svcGen.CreateServiceFile(generateFileBasicData, req.Table, req.Fields); err != nil {
		failLog("生成服务文件失败: " + err.Error())
		return
	}
	if err := svcGen.CreateHandlerFile(generateFileBasicData, req.Table, req.Fields); err != nil {
		failLog("生成控制器文件失败: " + err.Error())
		return
	}
	if err := svcGen.CreateRouterFile(generateFileBasicData, req.Table, req.Fields); err != nil {
		failLog("生成路由文件失败: " + err.Error())
		return
	}

	// 生成前端中英文语言包文件
	if err := svcGen.CreateLangFile(c.Request.Context(), generateFileBasicData, req.Table, req.Fields); err != nil {
		failLog("生成语言包文件失败: " + err.Error())
		return
	}

	// 生成 index.vue
	if err := svcGen.CreateIndexVueFile(generateFileBasicData, req.Table, req.Fields); err != nil {
		failLog("生成 index.vue 文件失败: " + err.Error())
		return
	}

	// 生成 dialogForm.vue
	if err := svcGen.CreateDialogFormFile(generateFileBasicData, req.Table, req.Fields); err != nil {
		failLog("生成 dialogForm.vue 文件失败: " + err.Error())
		return
	}

	// 同时格式化 index.vue / dialogForm.vue 的代码
	if err := filesystem.FormatWithPrettier(generateFileBasicData["views"].Dir); err != nil {
		failLog("格式化 index.vue / dialogForm.vue 文件失败: " + err.Error())
		return
	}

	// 写入后台菜单与权限节点
	if generateFileBasicData["router"].App == "admin" {
		if err := svcCrud.CreateMenuRule(c.Request.Context(), h.repoAuth, generateFileBasicData, req.Table, req.Fields); err != nil {
			failLog("写入菜单数据失败: " + err.Error())
			return
		}
	}

	// 标记生成日志为成功
	if logID > 0 {
		if _, err := svcCrud.CrudLog(c.Request.Context(), h.repo, dto.CrudLogData{
			ID:     &logID,
			Status: "succeeded",
		}); err != nil {
			failLog("更新生成日志失败: " + err.Error())
			return
		}
	}

	// 是否等待热更新
	isAir := false
	if pid, ok := airx.AirPID(); ok && pid > 0 {
		isAir = true
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"air":                   isAir,
		"req":                   req,
		"generateFileBasicData": generateFileBasicData,
	}))
}

// ChatStream AI 对话流式输出
// 转发上游 OpenAI Responses 兼容接口的 SSE 事件流给前端
func (h *Handler) ChatStream(c *gin.Context) {
	if !h.checkPermission(c, []string{"crud/crud/create"}) {
		return
	}
	var req dto.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	// SSE 流式响应头: 正常输出与错误事件均以 SSE 形式下发
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	body, err := h.svc.Stream(c.Request.Context(), &req)
	if err != nil {
		payload, _ := json.Marshal(map[string]any{
			"type":    "error",
			"message": err.Error(),
		})
		_, _ = c.Writer.WriteString("event: error\n")
		_, _ = c.Writer.WriteString("data: " + string(payload) + "\n\n")
		c.Writer.Flush()
		return
	}
	defer body.Close()

	// 默认 ScanLines 按行切分，正好对应 SSE 的 event/data 行（\n 与 \r\n 均兼容）
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			// 事件间空行，转发以分隔事件
			if _, err := c.Writer.WriteString("\n"); err != nil {
				return
			}
			c.Writer.Flush()
			continue
		}
		// 仅转发 event/data 行，其余（如注释行）忽略
		if strings.HasPrefix(line, "event:") || strings.HasPrefix(line, "data:") {
			if _, err := c.Writer.WriteString(line + "\n"); err != nil {
				return
			}
			c.Writer.Flush()
		}
	}
	if err := scanner.Err(); err != nil {
		// 上游流读取出错（如连接中断），流式输出到此结束
		return
	}
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
