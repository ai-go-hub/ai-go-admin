package routine

import (
	"strconv"

	"github.com/ai-go-hub/ai-go-admin/internal/handler"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	repoCommon "github.com/ai-go-hub/ai-go-admin/internal/repository/common"
	"github.com/ai-go-hub/ai-go-admin/internal/service"
	svcRoutine "github.com/ai-go-hub/ai-go-admin/internal/service/admin/routine"
	"github.com/ai-go-hub/ai-go-admin/pkg/jsonx"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

// ConfigHandler 系统配置控制器
type ConfigHandler struct {
	*handler.Handler[model.Config]
	svc  *svcRoutine.ConfigService
	repo *repoCommon.ConfigRepository
}

// NewConfigHandler 创建系统配置控制器实例
func NewConfigHandler(svc *svcRoutine.ConfigService, repo *repoCommon.ConfigRepository) *ConfigHandler {
	return &ConfigHandler{
		Handler: handler.NewHandler(svc),
		svc:     svc,
		repo:    repo,
	}
}

// RegisterRoutes 注册路由
func (h *ConfigHandler) RegisterRoutes(group *gin.RouterGroup) {
	handler.RegisterBaseRoutes(h, group)

	group.POST("/send-test-mail", middleware.AdminAuth(), middleware.AdminPermission(), h.SendTestMail)
}

// List 重写: 返回分组配置列表
func (h *ConfigHandler) List(c *gin.Context) {
	configs, err := h.svc.List(c.Request.Context(), service.Options{
		Limit:     999999,
		SortField: "weigh",
		SortOrder: "desc",
	})
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("查询配置失败: "+err.Error()))
		return
	}

	// 获取分组定义
	allGroups, err := h.repo.GetConfigGroups(c.Request.Context())
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("查询配置分组失败: "+err.Error()))
		return
	}

	// 返回数据组装为前端需要的格式
	var quickEntrance []map[string]string
	groupItems := make(map[string][]model.Config)

	for _, cfg := range configs {
		if cfg.Name == "quick_entrance" {
			quickEntrance = jsonx.UnmarshalSafe[[]map[string]string]([]byte(util.FromPtr(cfg.Value)))
		}
		groupItems[cfg.Group] = append(groupItems[cfg.Group], cfg)
	}

	var list []any
	for _, g := range allGroups {
		list = append(list, gin.H{
			"name":    g.Key,
			"title":   g.Value,
			"configs": groupItems[g.Key],
		})
	}

	configGroupMap := make(map[string]string)
	for _, g := range allGroups {
		configGroupMap[g.Key] = g.Value
	}

	if quickEntrance == nil {
		quickEntrance = []map[string]string{}
	}

	httpx.Success(c, httpx.WithData(gin.H{
		"list":          list,
		"configGroup":   configGroupMap,
		"quickEntrance": quickEntrance,
	}))
}

// Update 重写: 按 group 批量更新配置值
func (h *ConfigHandler) Update(c *gin.Context) {
	group := c.Param("pk")

	var entity map[string]string
	if err := c.ShouldBindJSON(&entity); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	if err := h.svc.BatchSave(c.Request.Context(), group, entity); err != nil {
		httpx.Fail(c, httpx.WithMessage("更新失败: "+err.Error()))
		return
	}

	httpx.Success(c, httpx.WithMessage("当前页面配置项更新成功"))
}

// SendTestMail 发送测试邮件
func (h *ConfigHandler) SendTestMail(c *gin.Context) {
	var req struct {
		SMTPServer     string `json:"smtp_server"`
		SMTPPort       string `json:"smtp_port"`
		SMTPUser       string `json:"smtp_user"`
		SMTPPass       string `json:"smtp_pass"`
		SMTPSenderMail string `json:"smtp_sender_mail"`
		SMTPVerify     string `json:"smtp_verification"`
		TestMail       string `json:"test_mail"`
		Name           string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, httpx.WithMessage("参数错误: "+err.Error()))
		return
	}

	port, err := strconv.Atoi(req.SMTPPort)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("SMTP 端口格式错误"))
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", req.SMTPSenderMail)
	m.SetHeader("To", req.TestMail)
	m.SetHeader("Subject", "这是一封测试邮件 - "+req.Name)
	m.SetBody("text/plain", "恭喜您，收到此邮件即表示您的邮件服务配置正确！")

	d := gomail.NewDialer(req.SMTPServer, port, req.SMTPUser, req.SMTPPass)
	if req.SMTPVerify == "SSL" {
		d.SSL = true
	}

	if err := d.DialAndSend(m); err != nil {
		httpx.Fail(c, httpx.WithMessage("邮件发送失败: "+err.Error()))
		return
	}
	httpx.Success(c, httpx.WithMessage("测试邮件发送成功~"))
}
