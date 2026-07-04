package database

import (
	"fmt"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var db *gorm.DB

// 初始化数据库连接
func Init() error {
	cfg := config.Get().Database

	gormCfg := &gorm.Config{}
	if cfg.Prefix != "" {
		gormCfg.NamingStrategy = schema.NamingStrategy{
			TablePrefix: cfg.Prefix,
		}
	}

	var err error
	db, err = gorm.Open(postgres.Open(buildDSN(cfg.Write)), gormCfg)
	if err != nil {
		return fmt.Errorf("connect write database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}

	// 连接池配置
	sqlDB.SetMaxOpenConns(cfg.Write.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Write.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Write.ConnMaxLifetime) * time.Second)

	// 配置读写分离
	if cfg.Read.Enabled {
		resolver := dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{postgres.Open(buildDSN(cfg.Read.DatabaseSource))},
		})
		resolver.SetMaxOpenConns(cfg.Read.MaxOpenConns)
		resolver.SetMaxIdleConns(cfg.Read.MaxIdleConns)
		resolver.SetConnMaxLifetime(time.Duration(cfg.Read.ConnMaxLifetime) * time.Second)

		if err := db.Use(resolver); err != nil {
			return fmt.Errorf("configure read database: %w", err)
		}
	}

	// 自动迁移所有已注册的模型
	if models := model.All(); len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			return fmt.Errorf("auto migrate: %w", err)
		}
	}

	return nil
}

// 返回全局数据库实例
func DB() *gorm.DB {
	return db
}

// 将 *gorm.DB 注入 gin.Context，绑定当前请求的 context
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", DB().WithContext(c.Request.Context()))
		c.Next()
	}
}

// 从 gin.Context 中获取带请求上下文的 *gorm.DB
func FromContext(c *gin.Context) *gorm.DB {
	return c.MustGet("db").(*gorm.DB)
}

// 根据配置构建 PostgreSQL DSN
func buildDSN(src config.DatabaseSource) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		src.Host, src.User, src.Password, src.DBName, src.Port, src.SSLMode, src.Timezone,
	)
}
