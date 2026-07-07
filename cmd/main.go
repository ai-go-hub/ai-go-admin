package main

import (
	"fmt"
	"os"

	"github.com/ai-go-hub/ai-go-admin/cmd/migrate"
	"github.com/ai-go-hub/ai-go-admin/cmd/serve"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:           "aigo",
		Short:         "命令行统一入口",
		Long:          "默认启动 API 服务，通过 migrate 子命令管理数据库迁移",
		SilenceUsage:  true,
		SilenceErrors: true,
		// 无子命令时默认启动 API 服务
		RunE: serve.Run,
	}
	root.AddCommand(serve.NewCommand())
	root.AddCommand(migrate.NewCommand())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "错误: ", err)
		os.Exit(1)
	}
}
