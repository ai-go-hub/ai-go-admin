package driver

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Local 本地磁盘上传驱动
type Local struct {
	dir       string // 磁盘存储根目录
	urlPrefix string // 访问 URL 前缀
}

// NewLocal 创建本地磁盘驱动
func NewLocal(dir, urlPrefix string) *Local {
	return &Local{dir: dir, urlPrefix: urlPrefix}
}

// Save 保存文件，storedFilename 为 / 分隔的相对路径，所在目录不存在时自动创建
func (l *Local) Save(reader io.Reader, storedFilename string) error {
	fullPath, err := l.FullPath(storedFilename)
	if err != nil {
		return err
	}

	if err := l.validatePath(fullPath); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("创建目录: %w", err)
	}

	// O_EXCL 防竞争条件: 文件或符号链接已存在则失败，不跟随链接写入
	dst, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("创建文件: %w", err)
	}
	defer dst.Close()
	if _, err := io.Copy(dst, reader); err != nil {
		return fmt.Errorf("写入文件: %w", err)
	}
	return nil
}

// Delete 删除文件，文件不存在视为成功
func (l *Local) Delete(storedFilename string) error {
	fullPath, err := l.FullPath(storedFilename)
	if err != nil {
		return err
	}

	// 删除前校验路径归属并拒绝符号链接
	if err := l.validatePath(fullPath); err != nil {
		return err
	}

	if _, err := os.Lstat(fullPath); err != nil && os.IsNotExist(err) {
		return nil
	}
	return os.Remove(fullPath)
}

// Url 返回文件的访问地址
func (l *Local) Url(storedFilename string) string {
	return l.urlPrefix + "/" + strings.TrimPrefix(storedFilename, "/")
}

// Exists 判断文件是否存在，路径不合法或命中符号链接时视为不存在
func (l *Local) Exists(storedFilename string) bool {
	fullPath, err := l.FullPath(storedFilename)
	if err != nil {
		return false
	}

	if err := l.validatePath(fullPath); err != nil {
		return false
	}

	_, err = os.Lstat(fullPath)
	return err == nil
}

// FullPath 返回文件在磁盘上的完整存储路径
func (l *Local) FullPath(path string) (string, error) {
	clean := strings.TrimPrefix(strings.TrimPrefix(path, l.urlPrefix), "/")
	if clean == "" || clean == "." {
		return "", fmt.Errorf("文件路径不能为空")
	}

	nativePath := filepath.FromSlash(clean)
	if filepath.IsAbs(nativePath) {
		return "", fmt.Errorf("文件路径不能是绝对路径: %s", path)
	}

	root, err := l.root()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, nativePath), nil
}

// validatePath 校验路径位于上传根目录内，且路径中不含符号链接
func (l *Local) validatePath(fullPath string) error {
	root, err := l.root()
	if err != nil {
		return err
	}

	// 确保路径位于根目录内，且不能是根目录本身
	rel, err := filepath.Rel(root, fullPath)
	if err != nil {
		return fmt.Errorf("校验路径归属: %w", err)
	}
	if rel == "." || rel == ".." || filepath.IsAbs(rel) || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("文件路径越出存储根目录: %s", fullPath)
	}

	// 逐层向上检查，拒绝路径中已存在的符号链接，避免通过链接访问根目录外的文件
	for current := fullPath; current != root; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("读取路径信息: %w", err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("文件路径不允许包含符号链接: %s", current)
		}
	}
	return nil
}

// root 返回上传根目录的绝对路径
func (l *Local) root() (string, error) {
	root, err := filepath.Abs(l.dir)
	if err != nil {
		return "", fmt.Errorf("解析上传根目录: %w", err)
	}
	return filepath.Clean(root), nil
}
