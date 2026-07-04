package common

import (
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"
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
