package crud

import (
	"context"
	"errors"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"

	"gorm.io/gorm"
)

// CrudLogRepository CRUD 记录仓储，嵌入通用仓储
type CrudLogRepository struct {
	*repository.Repository[model.CrudLog]
}

// NewCrudLogRepository 创建 CRUD 记录仓储实例
func NewCrudLogRepository() *CrudLogRepository {
	return &CrudLogRepository{
		Repository: repository.NewRepository[model.CrudLog](),
	}
}

// FindSucceededByName 根据表名查询生成成功的 CRUD 记录，未找到返回 nil
func (r *CrudLogRepository) FindSucceededByName(ctx context.Context, table string) (*model.CrudLog, error) {
	log, err := gorm.G[model.CrudLog](r.DB()).
		Where("name = ?", table).
		Where("status = ?", "succeeded").
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}
