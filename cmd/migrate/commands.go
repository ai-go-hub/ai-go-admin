package migrate

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

// upCmd 应用待执行的迁移
func upCmd() *cobra.Command {
	var steps int
	cmd := &cobra.Command{
		Use:   "up",
		Short: "应用数据库迁移（默认应用全部待执行迁移）",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := newMigrate()
			if err != nil {
				return err
			}
			defer m.Close()

			if steps > 0 {
				if err := m.Steps(steps); err != nil && err != migrate.ErrNoChange {
					return err
				}
			} else {
				if err := m.Up(); err != nil && err != migrate.ErrNoChange {
					return err
				}
			}
			fmt.Println("迁移完成。")
			return printVersion(m)
		},
	}
	cmd.Flags().IntVarP(&steps, "steps", "n", 0, "仅应用 N 个迁移（默认全部）")
	return cmd
}

// downCmd 回滚迁移
func downCmd() *cobra.Command {
	var all bool
	var steps int
	cmd := &cobra.Command{
		Use:   "down",
		Short: "回滚数据库迁移（默认回滚 1 个）",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := newMigrate()
			if err != nil {
				return err
			}
			defer m.Close()

			if all {
				if err := m.Down(); err != nil && err != migrate.ErrNoChange {
					return err
				}
			} else {
				n := steps
				if n <= 0 {
					n = 1
				}
				if err := m.Steps(-n); err != nil && err != migrate.ErrNoChange {
					return err
				}
			}
			fmt.Println("迁移回滚完成。")
			return printVersion(m)
		},
	}
	cmd.Flags().IntVarP(&steps, "steps", "n", 1, "回滚 N 个迁移（与 --all 互斥）")
	cmd.Flags().BoolVar(&all, "all", false, "回滚全部迁移")
	cmd.MarkFlagsMutuallyExclusive("all", "steps")
	return cmd
}

// migEntry 单个迁移文件的信息
type migEntry struct {
	version int
	name    string
}

// parseMigrations 从内嵌迁移文件中解析并按版本号排序
func parseMigrations() ([]migEntry, error) {
	entries, err := fs.ReadDir(migrationFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("读取内嵌迁移文件: %w", err)
	}

	var migs []migEntry
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		parts := strings.SplitN(name, "_", 2)
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		label := strings.TrimSuffix(strings.TrimPrefix(name, parts[0]+"_"), ".up.sql")
		migs = append(migs, migEntry{n, label})
	}
	sort.Slice(migs, func(i, j int) bool { return migs[i].version < migs[j].version })
	return migs, nil
}

// statusCmd 查看迁移状态
func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "查看迁移状态（当前版本与已应用/待执行列表）",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := newMigrate()
			if err != nil {
				return err
			}
			defer m.Close()

			if err := printVersion(m); err != nil {
				return err
			}
			fmt.Println()

			migs, err := parseMigrations()
			if err != nil {
				return err
			}

			v, _, err := m.Version()
			current := 0
			if err == nil {
				current = int(v)
			} else if err != migrate.ErrNilVersion {
				return err
			}

			fmt.Printf("%-8s  %-9s  %s\n", "版本", "状态", "名称")
			for _, mg := range migs {
				status := "待迁移"
				if current > 0 && mg.version <= current {
					status = "已迁移"
				}
				fmt.Printf("%-8d  %-9s  %s\n", mg.version, status, mg.name)
			}
			return nil
		},
	}
}

// forceCmd 强制设置迁移版本，用于修复脏数据状态
func forceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "force <version>",
		Short: "强制设置迁移版本（修复脏数据状态），-1 表示无版本",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("无效的版本号: %s", args[0])
			}
			m, err := newMigrate()
			if err != nil {
				return err
			}
			defer m.Close()

			if err := m.Force(version); err != nil {
				return err
			}
			fmt.Printf("已强制设置迁移版本为 %d。\n", version)
			return nil
		},
	}
	return cmd
}

// createCmd 在 migrations 目录创建一对新的 up/down 迁移 SQL 文件
func createCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "创建一对新的迁移 SQL 文件（版本号自动递增）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			dir := filepath.Join("cmd", "migrate", "migrations")

			// 找到下一个版本号
			migs, err := parseMigrations()
			if err != nil {
				return err
			}
			next := 1
			if len(migs) > 0 {
				next = migs[len(migs)-1].version + 1
			}
			ver := fmt.Sprintf("%06d", next)

			upFile := filepath.Join(dir, fmt.Sprintf("%s_%s.up.sql", ver, name))
			downFile := filepath.Join(dir, fmt.Sprintf("%s_%s.down.sql", ver, name))

			upContent := fmt.Sprintf("-- 创建 %s 相关表\n", name)
			downContent := fmt.Sprintf("-- 回滚: 删除 %s 相关表\n", name)

			if err := os.WriteFile(upFile, []byte(upContent), 0644); err != nil {
				return fmt.Errorf("创建 %s: %w", upFile, err)
			}
			if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
				return fmt.Errorf("创建 %s: %w", downFile, err)
			}

			fmt.Printf("已创建迁移文件:\n  %s\n  %s\n", upFile, downFile)
			return nil
		},
	}
}

// versionCmd 查看当前迁移版本
func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "查看当前迁移版本",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := newMigrate()
			if err != nil {
				return err
			}
			defer m.Close()

			return printVersion(m)
		},
	}
}

// dropCmd 删除数据库中所有表（危险操作，需谨慎）
func dropCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "drop",
		Short: "删除数据库中所有表（危险操作）",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := newMigrate()
			if err != nil {
				return err
			}
			defer m.Close()

			if err := m.Drop(); err != nil {
				return err
			}
			fmt.Println("已删除数据库中所有表。")
			return nil
		},
	}
}
