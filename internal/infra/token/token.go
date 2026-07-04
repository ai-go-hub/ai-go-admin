package token

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/token/driver"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
)

// 令牌类型标识
const (
	TypeAdminLogin = "admin_login" // 管理员登录令牌
)

// Driver 令牌存储驱动接口
type Driver interface {
	Create(ctx context.Context, token *model.Token) error
	Get(ctx context.Context, token string) (*model.Token, error)
	Delete(ctx context.Context, token string) error
	Clear(ctx context.Context, userID uint, tokenType string) error
	ClearExpired(ctx context.Context) error
}

// Service 令牌服务
type Service struct {
	driver Driver
}

// NewService 创建令牌服务
func NewService(driver Driver) *Service {
	return &Service{driver: driver}
}

// Create 创建令牌，入库前自动对 Token 做 SHA256
func (m *Service) Create(ctx context.Context, token *model.Token) error {
	// 清理过期令牌（Create 为低频操作，使用独立 context 不受请求生命周期影响）
	_ = m.driver.ClearExpired(context.Background())

	token.Token = sha256Hex(token.Token)
	return m.driver.Create(ctx, token)
}

// Get 获取令牌信息
func (m *Service) Get(ctx context.Context, token string) (*model.Token, error) {
	return m.driver.Get(ctx, sha256Hex(token))
}

// Check 检查令牌是否存在且未过期
func (m *Service) Check(ctx context.Context, token string) bool {
	t, err := m.Get(ctx, token)
	if err != nil || t == nil {
		return false
	}
	return time.Now().Before(t.ExpiredAt)
}

// Delete 删除令牌
func (m *Service) Delete(ctx context.Context, token string) error {
	return m.driver.Delete(ctx, sha256Hex(token))
}

// Clear 清除指定用户指定类型的所有令牌
func (m *Service) Clear(ctx context.Context, userID uint, tokenType string) error {
	return m.driver.Clear(ctx, userID, tokenType)
}

// sha256Hex 返回 raw 的 SHA256 十六进制字符串
func sha256Hex(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", sum)
}

// ==================== 全局单例 ====================

var (
	instance *Service
	once     sync.Once
)

// Manager 返回全局令牌管理器实例，首次调用时根据配置自动初始化
func Manager() *Service {
	once.Do(func() {
		instance = NewService(newDriver(config.Get().Token.Driver))
	})
	return instance
}

// newDriver 根据配置创建存储驱动
func newDriver(name string) Driver {
	switch name {
	default:
		return driver.NewDatabase()
	}
}
