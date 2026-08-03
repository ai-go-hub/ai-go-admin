package common

import (
	"context"
	"errors"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"
	"github.com/ai-go-hub/ai-go-admin/pkg/jsonx"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"

	"gorm.io/gorm"
)

// ConfigRepository 系统配置仓储
type ConfigRepository struct {
	*repository.Repository[model.Config]
}

// NewConfigRepository 创建系统配置仓储实例
func NewConfigRepository() *ConfigRepository {
	return &ConfigRepository{
		Repository: repository.NewRepository[model.Config](),
	}
}

// ConfigGroupItem 配置分组项
type ConfigGroupItem struct {
	Key   string
	Value string
}

// GetConfigGroups 获取配置分组定义（解析 config_group 配置项的 JSON 值）
func (r *ConfigRepository) GetConfigGroups(ctx context.Context) ([]ConfigGroupItem, error) {
	cfg, err := gorm.G[model.Config](r.DB()).
		Where("name = ?", "config_group").
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	value := util.FromPtr(cfg.Value)

	if value == "" {
		return nil, nil
	}
	items := jsonx.UnmarshalSafe[[]ConfigGroupItem]([]byte(value))
	return items, nil
}
