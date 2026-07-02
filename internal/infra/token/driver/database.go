package driver

import (
	"context"
	"errors"
	"time"

	"ai-go-mall/internal/infra/database"
	"ai-go-mall/internal/model"

	"gorm.io/gorm"
)

// Database 基于关系型数据库的令牌驱动
type Database struct{}

// NewDatabase 创建数据库令牌驱动
func NewDatabase() *Database {
	return &Database{}
}

// Create 创建令牌
func (d *Database) Create(ctx context.Context, t *model.Token) error {
	return gorm.G[model.Token](database.DB()).Create(ctx, t)
}

// Get 获取令牌信息
func (d *Database) Get(ctx context.Context, token string) (*model.Token, error) {
	t, err := gorm.G[model.Token](database.DB()).
		Where("token = ?", token).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// Delete 删除令牌
func (d *Database) Delete(ctx context.Context, token string) error {
	_, err := gorm.G[model.Token](database.DB()).Where("token = ?", token).Delete(ctx)
	return err
}

// Clear 清除指定用户指定类型的所有令牌
func (d *Database) Clear(ctx context.Context, userID uint, tokenType string) error {
	_, err := gorm.G[model.Token](database.DB()).
		Where("user_id = ? AND type = ?", userID, tokenType).
		Delete(ctx)
	return err
}

// ClearExpired 清除所有已过期的令牌
func (d *Database) ClearExpired(ctx context.Context) error {
	_, err := gorm.G[model.Token](database.DB()).
		Where("expired_at < ?", time.Now()).
		Delete(ctx)
	return err
}
