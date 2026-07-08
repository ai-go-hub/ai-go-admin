package migrate

import (
	"embed"
	"fmt"
	"net/url"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/cobra"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// NewCommand 返回数据库迁移命令组
func NewCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "migrate",
		Short:         "数据库迁移工具，基于 golang-migrate",
		Long:          "基于 golang-migrate 的数据库迁移工具，支持 up / down / status / version / drop / force / create 子命令",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Init(); err != nil {
				return fmt.Errorf("初始化配置: %w", err)
			}
			loc, err := time.LoadLocation(config.Get().App.Timezone)
			if err != nil {
				return fmt.Errorf("加载时区: %w", err)
			}
			time.Local = loc
			return nil
		},
	}
	root.AddCommand(upCmd(), downCmd(), statusCmd(), versionCmd(), dropCmd(), forceCmd(), createCmd())
	return root
}

// newMigrate 基于内嵌迁移文件与项目数据库配置构造迁移实例
func newMigrate() (*migrate.Migrate, error) {
	src, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("加载内嵌迁移文件: %w", err)
	}

	// 将 __PREFIX__ 占位符替换为数据库表前缀
	src = &prefixDriver{
		Driver: src,
		prefix: config.Get().Database.Prefix,
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dsn())
	if err != nil {
		return nil, fmt.Errorf("连接数据库: %w", err)
	}
	return m, nil
}

// dsn 依据项目数据库写库配置生成 PostgreSQL 连接 URL
func dsn() string {
	c := config.Get().Database.Write
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   c.DBName,
	}
	q := u.Query()
	q.Set("sslmode", c.SSLMode)
	if tz := config.Get().App.Timezone; tz != "" {
		q.Set("TimeZone", tz)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// printVersion 打印当前迁移版本
func printVersion(m *migrate.Migrate) error {
	v, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			fmt.Println("当前版本: 0")
			return nil
		}
		return err
	}
	state := "干净（无迁移失败）"
	if dirty {
		state = "存在脏数据（有迁移失败）"
	}
	fmt.Printf("当前版本: %d - %s\n", v, state)
	return nil
}
