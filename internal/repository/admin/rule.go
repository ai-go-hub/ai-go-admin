package admin

import (
	"context"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"

	"gorm.io/gorm"
)

// AdminRuleRepository 菜单和权限规则仓储，嵌入通用仓储
type AdminRuleRepository struct {
	*repository.Repository[model.AdminRule]
}

// NewAdminRuleRepository 创建菜单和权限规则仓储实例
func NewAdminRuleRepository() *AdminRuleRepository {
	return &AdminRuleRepository{
		Repository: repository.NewRepository[model.AdminRule](),
	}
}

// FindByName 根据规则名称查询
func (r *AdminRuleRepository) FindByName(ctx context.Context, name string) (*model.AdminRule, error) {
	rule, err := gorm.G[model.AdminRule](r.DB()).Where("name = ?", name).First(ctx)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// FindByPath 根据菜单路由路径查询
func (r *AdminRuleRepository) FindByPath(ctx context.Context, path string) (*model.AdminRule, error) {
	rule, err := gorm.G[model.AdminRule](r.DB()).Where("path = ?", path).First(ctx)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// DistinctPidsByIDs 返回 ids 的所有 pid 的去重列表
// 用于父子关系剥离等场景，NULL 的 pid 会被过滤
func (r *AdminRuleRepository) DistinctPidsByIDs(ctx context.Context, ids []uint) ([]uint, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var pids []uint
	err := r.DB().WithContext(ctx).
		Model(&model.AdminRule{}).
		Where("id IN ? AND pid IS NOT NULL", ids).
		Distinct().
		Pluck("pid", &pids).Error
	if err != nil {
		return nil, err
	}
	return pids, nil
}

// ChildIDsByPids 根据 pids 查直接子级 id 集合
func (r *AdminRuleRepository) ChildIDsByPids(ctx context.Context, pids []uint) ([]uint, error) {
	if len(pids) == 0 {
		return nil, nil
	}
	children, err := gorm.G[model.AdminRule](r.DB()).
		Where("pid IN ?", pids).
		Select("id").
		Find(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(children))
	for _, rule := range children {
		ids = append(ids, rule.ID)
	}
	return ids, nil
}

// TitleMapByIDs 返回 ids 中指定 types 的规则的 id -> title 映射
func (r *AdminRuleRepository) TitleMapByIDs(ctx context.Context, ids []uint, types ...string) (map[uint]string, error) {
	if len(ids) == 0 {
		return map[uint]string{}, nil
	}
	q := gorm.G[model.AdminRule](r.DB()).Where("id IN ?", ids)
	if len(types) > 0 {
		q = q.Where("type IN ?", types)
	}
	rules, err := q.Find(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[uint]string, len(rules))
	for _, rule := range rules {
		result[rule.ID] = rule.Title
	}
	return result, nil
}
