package admin

import (
	"context"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"

	"gorm.io/gorm"
)

// AdminGroupRepository 管理员分组仓储，嵌入通用仓储
type AdminGroupRepository struct {
	*repository.Repository[model.AdminGroup]
}

// NewAdminGroupRepository 创建管理员分组仓储实例
func NewAdminGroupRepository() *AdminGroupRepository {
	return &AdminGroupRepository{
		Repository: repository.NewRepository[model.AdminGroup](),
	}
}

// FindByName 根据组名查询分组
func (r *AdminGroupRepository) FindByName(ctx context.Context, name string) (*model.AdminGroup, error) {
	group, err := gorm.G[model.AdminGroup](r.DB()).Where("name = ?", name).First(ctx)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// FindByIDs 根据 ID 列表批量查询分组
func (r *AdminGroupRepository) FindByIDs(ctx context.Context, ids []uint) ([]model.AdminGroup, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return gorm.G[model.AdminGroup](r.DB()).Where("id IN ?", ids).Find(ctx)
}

// FindAll 查询全部分组（不带过滤条件），用于内存中构建父子关系
func (r *AdminGroupRepository) FindAll(ctx context.Context) ([]model.AdminGroup, error) {
	return gorm.G[model.AdminGroup](r.DB()).Find(ctx)
}

// ChildIDsByPids 根据 pids 查直接子级 id 集合
func (r *AdminGroupRepository) ChildIDsByPids(ctx context.Context, pids []uint) ([]uint, error) {
	if len(pids) == 0 {
		return nil, nil
	}
	children, err := gorm.G[model.AdminGroup](r.DB()).
		Where("pid IN ?", pids).
		Select("id").
		Find(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(children))
	for _, g := range children {
		ids = append(ids, g.ID)
	}
	return ids, nil
}
