package common

import (
	"fmt"
	"hash/fnv"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/kit/httpx"
	repoCommon "github.com/ai-go-hub/ai-go-admin/internal/repository/common"
	"github.com/ai-go-hub/ai-go-admin/pkg/xss"
	"github.com/gin-gonic/gin"
)

// FileSvg 根据文件后缀生成 SVG 文件图标
func FileSvg(c *gin.Context) {
	suffix := c.DefaultQuery("suffix", "file")
	background := c.Query("background")

	c.Header("Content-Type", "image/svg+xml; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=604800")
	c.String(http.StatusOK, buildSuffixSvg(suffix, background))
}

// UtilHandler 工具控制器
type UtilHandler struct {
	configRepo *repoCommon.ConfigRepository
}

// NewUtilHandler 创建工具控制器实例
func NewUtilHandler() *UtilHandler {
	return &UtilHandler{
		configRepo: repoCommon.NewConfigRepository(),
	}
}

// RegisterRoutes 注册工具路由
func (h *UtilHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/area", Area)
	group.GET("/file-svg", FileSvg)
	group.GET("/site-config", h.SiteConfig)
}

// SiteConfig 获取站点配置数据
func (h *UtilHandler) SiteConfig(c *gin.Context) {
	siteConfigNames := []string{"name", "record_number", "ps_record_number", "version"}
	result, err := h.configRepo.GetConfigs(c.Request.Context(), siteConfigNames)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("获取站点配置失败: "+err.Error()))
		return
	}

	result["timezone"] = config.Get().App.Timezone
	result["cdn_url"] = config.Get().CDN.URL
	result["cdn_url_params"] = config.Get().CDN.URLParams

	httpx.Success(c, httpx.WithData(result))
}

// buildSuffixSvg 构建文件后缀的 SVG 图片
func buildSuffixSvg(suffix, background string) string {
	// 清理 XSS
	suffix = xss.TextPolicySanitize(suffix)
	background = xss.TextPolicySanitize(background)

	suffix = strings.ToUpper(suffix)
	if len(suffix) > 4 {
		suffix = suffix[:4]
	}

	// 基于后缀生成确定性色调
	h := fnv.New32a()
	h.Write([]byte(suffix))
	total := h.Sum32()
	hue := int(total % 360)

	r, g, b := hsvToRGB(float64(hue)/360.0, 0.3, 0.9)

	if background == "" {
		background = fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	}

	return fmt.Sprintf(`<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px" viewBox="0 0 512 512" style="enable-background:new 0 0 512 512;" xml:space="preserve">
            <path style="fill:#E2E5E7;" d="M128,0c-17.6,0-32,14.4-32,32v448c0,17.6,14.4,32,32,32h320c17.6,0,32-14.4,32-32V128L352,0H128z"/>
            <path style="fill:#B0B7BD;" d="M384,128h96L352,0v96C352,113.6,366.4,128,384,128z"/>
            <path style="fill:%s;" d="M416,416c0,8.8-7.2,16-16,16H48c-8.8,0-16-7.2-16-16V256c0-8.8,7.2-16,16-16h352c8.8,0,16,7.2,16,16 V416z"/>
            <g><text><tspan x="220" y="380" font-size="124" font-family="Verdana, Helvetica, Arial, sans-serif" fill="white" text-anchor="middle">%s</tspan></text></g>
        </svg>`, background, suffix)
}

// Area 获取省份地区数据
func Area(c *gin.Context) {
	city := c.DefaultQuery("city", "")
	province := c.DefaultQuery("province", "")

	pid := 0
	level := 1
	if province != "" {
		pid, _ = strconv.Atoi(province)
		level = 2
		if city != "" {
			pid, _ = strconv.Atoi(city)
			level = 3
		}
	}

	areas, err := repoCommon.NewAreaRepository().FindByPidAndLevel(c.Request.Context(), pid, level)
	if err != nil {
		httpx.Fail(c, httpx.WithMessage("地区数据查询失败"))
		return
	}

	httpx.Success(c, httpx.WithData(areas))
}

// hsvToRGB 将 HSV 颜色转换为 RGB 整数值
func hsvToRGB(h, s, v float64) (r, g, b int) {
	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)

	var rf, gf, bf float64
	switch int(i) % 6 {
	case 0:
		rf, gf, bf = v, t, p
	case 1:
		rf, gf, bf = q, v, p
	case 2:
		rf, gf, bf = p, v, t
	case 3:
		rf, gf, bf = p, q, v
	case 4:
		rf, gf, bf = t, p, v
	case 5:
		rf, gf, bf = v, p, q
	}

	return int(math.Floor(rf * 255)), int(math.Floor(gf * 255)), int(math.Floor(bf * 255))
}
