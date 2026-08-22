package routine

import (
	"context"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	repoCommon "github.com/ai-go-hub/ai-go-admin/internal/repository/common"
	"github.com/ai-go-hub/ai-go-admin/internal/service"
	"github.com/ai-go-hub/ai-go-admin/pkg/xss"

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
	// 预加载该分组配置的 name -> type 映射，对 editor 类型值做反 XSS 清洗
	list, err := gorm.G[model.Config](s.repo.DB()).
		Where(`"group" = ?`, group).
		Find(ctx)
	if err != nil {
		return err
	}
	types := make(map[string]string, len(list))
	for _, cfg := range list {
		types[cfg.Name] = cfg.Type
	}

	for name, value := range data {
		if types[name] == "editor" {
			value = xss.HTMLPolicySanitize(value)
		}
		_, err := gorm.G[map[string]any](s.repo.DB()).
			Table(model.Config{}.TableName()).
			Where(`name = ? AND "group" = ?`, name, group).
			Updates(ctx, map[string]any{"value": value})
		if err != nil {
			return err
		}
	}
	return nil
}
