package main

import (
	"fmt"
	"log"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"
	"github.com/ai-go-hub/ai-go-admin/internal/middleware"
	"github.com/ai-go-hub/ai-go-admin/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化文件配置系统
	if err := config.Init(); err != nil {
		log.Fatalf("config init: %v", err)
	}

	// 初始化数据库连接
	if err := database.Init(); err != nil {
		log.Fatalf("database init: %v", err)
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
		log.Fatalf("server start: %v", err)
	}
}
