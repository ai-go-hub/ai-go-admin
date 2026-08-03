package routine

import (
	"context"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	repoCommon "github.com/ai-go-hub/ai-go-admin/internal/repository/common"
	"github.com/ai-go-hub/ai-go-admin/internal/service"

	"gorm.io/gorm"
)

// ConfigService 系统配置服务
type ConfigService struct {
	service.IService[model.Config]
	repo *repoCommon.ConfigRepository
}

// NewConfigService 创建系统配置服务实例
func NewConfigService(repo *repoCommon.ConfigRepository) *ConfigService {
	return &ConfigService{
		IService: service.NewService(repo),
		repo:     repo,
	}
}

// BatchSave 批量保存配置值: 根据 group 和 name 定位行，更新 value
func (s *ConfigService) BatchSave(ctx context.Context, group string, data map[string]string) error {
	for name, value := range data {
		_, err := gorm.G[map[string]any](s.repo.DB()).
			Table(model.Config{}.TableName()).
			Where("name = ? AND \"group\" = ?", name, group).
			Updates(ctx, map[string]any{"value": value})
		if err != nil {
			return err
		}
	}
	return nil
}
