package serve

import (
	"fmt"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	"github.com/ai-go-hub/ai-go-admin/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// NewCommand 返回启动 API 服务的 serve 命令
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:           "serve",
		Short:         "启动 API 服务",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          Run,
	}
}

// Run 启动 API 服务
func Run(cmd *cobra.Command, args []string) error {
	// 初始化文件配置系统
	if err := config.Init(); err != nil {
		return fmt.Errorf("初始化配置: %w", err)
	}

	// 设置全局时区
	loc, err := time.LoadLocation(config.Get().App.Timezone)
	if err != nil {
		return fmt.Errorf("加载时区: %w", err)
	}
	time.Local = loc

	// 初始化数据库连接
	if err := database.Init(); err != nil {
		return fmt.Errorf("初始化数据库: %w", err)
	}

	engine := gin.Default()

	// 注册跨域中间件
	engine.Use(middleware.CORS())

	// 注册数据库中间件
	engine.Use(database.Middleware())

	// 注册路由
	router.Setup(engine)

	cfg := config.Get()
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	if err := engine.Run(addr); err != nil {
		return fmt.Errorf("启动服务: %w", err)
	}
	return nil
}
