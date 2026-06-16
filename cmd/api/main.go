package main

import (
	"fmt"
	"log"

	"ai-go-mall/config"
	"ai-go-mall/internal/infra/database"
	"ai-go-mall/internal/router"

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
