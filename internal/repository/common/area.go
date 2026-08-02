package common

import (
	"context"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"

	"gorm.io/gorm"
)

// AreaRepository 省份地区仓储，嵌入通用仓储并扩展地区专属查询
type AreaRepository struct {
	*repository.Repository[model.Area]
}

// NewAreaRepository 创建省份地区仓储实例
func NewAreaRepository() *AreaRepository {
	return &AreaRepository{
		Repository: repository.NewRepository[model.Area](),
	}
}

// FindByPidAndLevel 根据上级ID和等级查询地区列表
func (r *AreaRepository) FindByPidAndLevel(ctx context.Context, pid int, level int) ([]model.Area, error) {
	return gorm.G[model.Area](r.DB()).
		Where("pid = ? AND level = ?", pid, level).
		Find(ctx)
}
