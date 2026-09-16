package common

import (
	"context"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"

	"gorm.io/gorm"
)

// CaptchaRepository 验证码仓储
type CaptchaRepository struct {
	*repository.Repository[model.Captcha]
}

// NewCaptchaRepository 创建验证码仓储实例
func NewCaptchaRepository() *CaptchaRepository {
	return &CaptchaRepository{
		Repository: repository.NewRepository[model.Captcha](),
	}
}

// GetByKey 根据 key 查询验证码记录
func (r *CaptchaRepository) GetByKey(ctx context.Context, key string) (*model.Captcha, error) {
	record, err := gorm.G[model.Captcha](r.DB()).Where("key = ?", key).First(ctx)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// DeleteByKey 根据 key 删除验证码记录
func (r *CaptchaRepository) DeleteByKey(ctx context.Context, key string) error {
	_, err := gorm.G[model.Captcha](r.DB()).Where("key = ?", key).Delete(ctx)
	return err
}

// Save 保存验证码记录
func (r *CaptchaRepository) Save(ctx context.Context, record *model.Captcha) error {
	return gorm.G[model.Captcha](r.DB()).Create(ctx, record)
}

// DeleteExpired 删除已过期的验证码记录
func (r *CaptchaRepository) DeleteExpired(ctx context.Context) error {
	_, err := gorm.G[model.Captcha](r.DB()).Where("expired_at < ?", time.Now()).Delete(ctx)
	return err
}
