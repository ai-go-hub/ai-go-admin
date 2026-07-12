package filesystem

import (
	"os"
	"path/filepath"
	"strconv"
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

// FormatBytes 将字节数格式化为易读单位（B/KB/MB/GB/TB/PB）
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return strconv.FormatInt(bytes, 10) + " B"
	}
	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit && exp < len(units)-1; n /= unit {
		div *= unit
		exp++
	}
	// 保留 2 位小数并去掉末尾多余的 0 和小数点
	s := strconv.FormatFloat(float64(bytes)/float64(div), 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s + " " + units[exp]
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
