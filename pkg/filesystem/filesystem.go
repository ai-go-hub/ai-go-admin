package filesystem

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// TrimExt 返回去除了路径和扩展名的文件名
// path:文件路径或完整文件名
func TrimExt(path string) string {
	name := filepath.Base(path)
	ext := filepath.Ext(name)
	return name[:len(name)-len(ext)]
}

// Exists 判断文件/目录是否存在
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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

// AbsInProject 将路径规范化为项目根目录下的绝对路径
func AbsInProject(rel string) (string, bool) {
	abs, err := filepath.Abs(rel)
	if err != nil {
		return "", false
	}

	root, err := os.Getwd()
	if err != nil {
		return "", false
	}

	// 仅允许位于项目根目录内，防止越界
	if abs != root && !strings.HasPrefix(abs, root+string(filepath.Separator)) {
		return "", false
	}
	return abs, true
}

// RemoveFileWithRetry 删除文件并返回是否删除成功；Windows 下句柄可能被短暂占用，故重试数次
func RemoveFileWithRetry(path string) bool {
	for range 5 {
		if err := os.Remove(path); err == nil {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}

// RemoveDirWithRetry 删除空目录；Windows 下目录句柄可能被短暂占用，故重试数次
func RemoveDirWithRetry(dir string) bool {
	for range 5 {
		if err := os.Remove(dir); err == nil {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}

// RemoveAllWithRetry 删除 root 之下的所有子项；path 必须是 root 的严格子路径
func RemoveAllWithRetry(path, root string) bool {
	path, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	if root == "" {
		root, err = os.Getwd()
	} else {
		root, err = filepath.Abs(root)
	}
	if err != nil {
		return false
	}

	if path == root || !strings.HasPrefix(path, root+string(filepath.Separator)) {
		return false
	}
	for range 5 {
		if err := os.RemoveAll(path); err == nil {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}

// RemoveEmptyParents 从文件所在目录向上删除已为空的目录，遇非空目录或项目根目录时停止
func RemoveEmptyParents(path string) {
	root, err := os.Getwd()
	if err != nil {
		return
	}
	for dir := filepath.Dir(path); ; dir = filepath.Dir(dir) {
		rel, err := filepath.Rel(root, dir)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			return
		}
		if !RemoveDirWithRetry(dir) {
			return
		}
	}
}

// RemoveFileWithCleanup 删除项目内的文件，并清理因删除而变空的目录
func RemoveFileWithCleanup(rel string) {
	path, ok := AbsInProject(rel)
	if !ok {
		return
	}
	if !RemoveFileWithRetry(path) {
		return
	}
	RemoveEmptyParents(path)
}
