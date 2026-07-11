package upload

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/upload/driver"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	repoCommon "github.com/ai-go-hub/ai-go-admin/internal/repository/common"
	"github.com/ai-go-hub/ai-go-admin/pkg/filesystem"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"
)

// Driver 上传存储驱动接口
type Driver interface {
	// Save 保存文件，storedFilename 为 / 分隔的相对路径
	Save(reader io.Reader, storedFilename string) error
	// Delete 删除文件
	Delete(storedFilename string) error
	// Url 返回文件的访问地址
	Url(storedFilename string) string
	// Exists 判断文件是否存在
	Exists(storedFilename string) bool
	// FullPath 返回文件在磁盘上的完整存储路径
	FullPath(storedFilename string) string
}

// Result 上传结果
type Result struct {
	Url     string // 资源访问地址
	Size    int64  // 文件大小（字节）
	Suffix  string // 文件后缀
	IsImage bool   // 是否为图片
}

// Service 上传服务
type Service struct {
	repo *repoCommon.AttachmentRepository
}

// NewService 创建上传服务
func NewService(repo *repoCommon.AttachmentRepository) *Service {
	return &Service{repo: repo}
}

// Upload 上传文件，校验大小与后缀白名单，按 sha1 去重（命中则引用计数 +1 并复用已存文件），入库附件记录
// driverName 指定上传驱动名称（如 local），topic 为业务分类（如 avatar、article），userID/userType 标识上传者
func (s *Service) Upload(ctx context.Context, header *multipart.FileHeader, driverName string, topic string, userID uint, userType string) (*Result, error) {
	cfg := config.Get().Upload

	// 校验文件大小
	if cfg.MaxSize > 0 && header.Size > int64(cfg.MaxSize) {
		return nil, fmt.Errorf("文件大小 %d 超过最大上传限制 %d", header.Size, cfg.MaxSize)
	}

	// 校验后缀白名单
	ext := filesystem.Extension(header.Filename)
	if !s.AllowSuffix(ext, cfg.AllowSuffix) {
		return nil, fmt.Errorf("不允许的文件后缀: %s", ext)
	}

	// 读取文件内容（计算 SHA1 需要完整内容）
	src, err := header.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件: %w", err)
	}
	defer src.Close()
	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取上传文件: %w", err)
	}

	// 按 sha1 查询是否已存在（秒传去重）
	now := time.Now()
	sha1Hex := sha1HexOf(data)
	att, err := s.repo.FindBySha1(ctx, sha1Hex)
	if err != nil {
		return nil, fmt.Errorf("查询附件: %w", err)
	}
	if att != nil {
		// 已存在: 引用计数 +1，更新最后上传时间与上传者，复用已存文件
		if err := s.repo.IncrementQuote(ctx, att.ID, now, userID, userType); err != nil {
			return nil, fmt.Errorf("更新附件引用: %w", err)
		}
		return s.buildResult(att.URL, att.Size, ext), nil
	}

	// 新文件: 实例化驱动并保存
	d, err := newDriver(driverName)
	if err != nil {
		return nil, err
	}
	storedFilename := s.FormatFilename(cfg.Filename, header.Filename, topic, data)
	if err := d.Save(bytes.NewReader(data), storedFilename); err != nil {
		return nil, err
	}

	// 入库
	att = &model.Attachment{
		Topic:        topic,
		UserID:       userID,
		UserType:     userType,
		URL:          d.Url(storedFilename),
		Name:         header.Filename,
		Size:         header.Size,
		Mimetype:     header.Header.Get("Content-Type"),
		Quote:        1,
		Driver:       driverName,
		Sha1:         sha1Hex,
		CreateAt:     now,
		LastUploadAt: now,
	}
	if err := s.repo.Save(ctx, att); err != nil {
		return nil, fmt.Errorf("入库附件: %w", err)
	}

	return s.buildResult(att.URL, header.Size, ext), nil
}

// buildResult 构造上传结果
func (s *Service) buildResult(url string, size int64, ext string) *Result {
	return &Result{
		Url:     url,
		Size:    size,
		Suffix:  "." + ext,
		IsImage: filesystem.IsImageExtension(ext),
	}
}

// sha1HexOf 返回数据的 SHA1 十六进制字符串
func sha1HexOf(data []byte) string {
	sum := sha1.Sum(data)
	return hex.EncodeToString(sum[:])
}

// AllowSuffix 判断扩展名是否在白名单内，白名单为空时禁止任何类型上传
func (s *Service) AllowSuffix(ext string, allow []string) bool {
	if len(allow) == 0 {
		return false
	}
	for _, a := range allow {
		if strings.ToLower(a) == ext {
			return true
		}
	}
	return false
}

// FormatFilename 按配置模板生成存储文件名
func (s *Service) FormatFilename(tpl, filename, topic string, data []byte) string {
	now := time.Now()
	repl := map[string]string{
		// 业务的主题或分类
		"{topic}": sanitize(topic),
		"{year}":  fmt.Sprintf("%04d", now.Year()),
		"{mon}":   fmt.Sprintf("%02d", int(now.Month())),
		"{day}":   fmt.Sprintf("%02d", now.Day()),
		// 文件名称最多保留前 15 个字符（按 rune 截断，兼容中文）
		"{fileName}": util.TruncateString(sanitize(filesystem.TrimExt(filename)), 15),
		"{fileSha1}": sha1HexOf(data),
		"{.suffix}":  "." + filesystem.Extension(filename),
	}
	out := tpl
	for k, v := range repl {
		out = strings.ReplaceAll(out, k, v)
	}
	return path.Clean(out)
}

// unsafeRe 匹配文件名中不利于路径安全与 URL 传输的字符
var unsafeRe = regexp.MustCompile(`[/\\:@#?&=]|\.\.`)

// sanitize 清理文件名片段，避免路径穿越并保证 URL 安全
func sanitize(s string) string {
	return unsafeRe.ReplaceAllString(s, "_")
}

// ==================== 全局单例 ====================

var (
	instance *Service
	once     sync.Once
)

// Manager 返回全局上传服务实例
func Manager() *Service {
	once.Do(func() {
		instance = NewService(repoCommon.NewAttachmentRepository())
	})
	return instance
}

// newDriver 根据驱动名称创建上传驱动
func newDriver(name string) (Driver, error) {
	switch name {
	case "local":
		return driver.NewLocal(), nil
	default:
		return nil, fmt.Errorf("不支持的上传驱动: %s", name)
	}
}
