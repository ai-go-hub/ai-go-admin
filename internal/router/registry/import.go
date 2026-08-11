package registry

import (
	"os"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/pkg/filesystem"
)

// RouterFile 路由注册入口文件
const RouterFile = "internal/router/router.go"

// AddRouterImport 向路由注册入口添加空白导入（已存在时忽略），并用 go fmt 格式化
func AddRouterImport(importPath string) error {
	content, err := os.ReadFile(RouterFile)
	if err != nil {
		return err
	}
	line := `_ "` + importPath + `"`
	if strings.Contains(string(content), line) {
		return nil
	}

	// 插入 import 块（顺序由后续 go fmt 自动整理）
	lines := strings.Split(string(content), "\n")
	idx := len(lines)
	inImport := false
	for i, ln := range lines {
		trimmed := strings.TrimSpace(ln)
		if trimmed == "import (" {
			inImport = true
			continue
		}
		if !inImport {
			continue
		}
		if strings.HasPrefix(trimmed, `_ "`) || trimmed == ")" {
			idx = i
			break
		}
	}
	out := make([]string, 0, len(lines)+1)
	out = append(out, lines[:idx]...)
	out = append(out, "\t"+line)
	out = append(out, lines[idx:]...)
	if err := os.WriteFile(RouterFile, []byte(strings.Join(out, "\n")), 0o644); err != nil {
		return err
	}
	return filesystem.FormatGoFile(RouterFile)
}

// RemoveRouterImport 从路由注册入口移除空白导入，并用 go fmt 格式化
func RemoveRouterImport(importPath string) error {
	content, err := os.ReadFile(RouterFile)
	if err != nil {
		return err
	}
	line := `_ "` + importPath + `"`

	lines := strings.Split(string(content), "\n")
	out := make([]string, 0, len(lines))
	removed := false
	for _, ln := range lines {
		if strings.TrimSpace(ln) == line {
			removed = true
			continue
		}
		out = append(out, ln)
	}
	if !removed {
		return nil
	}
	if err := os.WriteFile(RouterFile, []byte(strings.Join(out, "\n")), 0o644); err != nil {
		return err
	}
	return filesystem.FormatGoFile(RouterFile)
}
