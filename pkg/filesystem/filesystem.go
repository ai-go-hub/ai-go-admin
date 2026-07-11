package filesystem

import (
	"os"
	"path/filepath"
	"strings"
)

// TrimExt 返回去除了路径和扩展名的文件名
// path:文件路径或完整文件名
func TrimExt(path string) string {
	name := filepath.Base(path)
	ext := filepath.Ext(name)
	return name[:len(name)-len(ext)]
}

// Extension 返回文件扩展名（不含点、小写）
func Extension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	return strings.TrimPrefix(ext, ".")
}

// IsImageExtension 判断扩展名是否为图片类型
func IsImageExtension(ext string) bool {
	switch strings.ToLower(ext) {
	case "jpg", "jpeg", "png", "gif", "webp", "bmp", "svg", "avif":
		return true
	}
	return false
}

// ReadDir 递归读取目录下所有文件，返回完整路径列表
func ReadDir(dir string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}
