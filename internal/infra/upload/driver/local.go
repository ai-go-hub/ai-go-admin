package driver

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// root 本地磁盘存储根目录
const root = "static"

// Local 本地磁盘上传驱动
type Local struct{}

// NewLocal 创建本地磁盘驱动
func NewLocal() *Local {
	return &Local{}
}

// Save 保存文件，storedFilename 为 / 分隔的相对路径，所在目录不存在时自动创建
func (l *Local) Save(reader io.Reader, storedFilename string) error {
	fullPath := l.FullPath(storedFilename)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("创建目录: %w", err)
	}
	dst, err := os.Create(fullPath)
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
	if err := os.Remove(l.FullPath(storedFilename)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Url 返回文件的访问地址，直接使用存储文件名
func (l *Local) Url(storedFilename string) string {
	return storedFilename
}

// Exists 判断文件是否存在
func (l *Local) Exists(storedFilename string) bool {
	_, err := os.Stat(l.FullPath(storedFilename))
	return err == nil
}

// FullPath 返回文件在磁盘上的完整存储路径
func (l *Local) FullPath(storedFilename string) string {
	return filepath.Join(root, filepath.FromSlash(storedFilename))
}
