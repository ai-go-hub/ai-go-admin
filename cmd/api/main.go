package main

import (
	"fmt"
	"log"

	"ai-go-mall/config"
	"ai-go-mall/internal/database"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("config init: %v", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("database init: %v", err)
	}

	cfg := config.Get()

	fmt.Printf("应用名称: %s\n", cfg.App.Name)
	fmt.Printf("服务端口: %d\n", cfg.Server.Port)
	fmt.Printf("数据库类型: %s\n", cfg.Database.Type)
	fmt.Printf("写库地址: %s:%d\n", cfg.Database.Write.Host, cfg.Database.Write.Port)
	fmt.Printf("写库名称: %s\n", cfg.Database.Write.DBName)
	fmt.Printf("写库连接池: max_open=%d max_idle=%d max_lifetime=%ds\n",
		cfg.Database.Write.MaxOpenConns,
		cfg.Database.Write.MaxIdleConns,
		cfg.Database.Write.ConnMaxLifetime,
	)

	if cfg.Database.Read.Enabled {
		fmt.Printf("读库已启用: %s:%d\n", cfg.Database.Read.Host, cfg.Database.Read.Port)
		fmt.Printf("读库连接池: max_open=%d max_idle=%d max_lifetime=%ds\n",
			cfg.Database.Read.MaxOpenConns,
			cfg.Database.Read.MaxIdleConns,
			cfg.Database.Read.ConnMaxLifetime,
		)
	} else {
		fmt.Println("读库未启用")
	}
}
