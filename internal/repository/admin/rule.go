package admin

import (
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"

	"github.com/gin-gonic/gin"
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
func (r *AdminRuleRepository) FindByName(c *gin.Context, name string) (*model.AdminRule, error) {
	rule, err := gorm.G[model.AdminRule](r.DB()).Where("name = ?", name).First(c.Request.Context())
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// FindByPath 根据菜单路由路径查询
func (r *AdminRuleRepository) FindByPath(c *gin.Context, path string) (*model.AdminRule, error) {
	rule, err := gorm.G[model.AdminRule](r.DB()).Where("path = ?", path).First(c.Request.Context())
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// DistinctPidsByIDs 返回 ids 的所有 pid 的去重列表
// 用于父子关系剥离等场景，NULL 的 pid 会被过滤
func (r *AdminRuleRepository) DistinctPidsByIDs(c *gin.Context, ids []uint) ([]uint, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var pids []uint
	err := r.DB().WithContext(c.Request.Context()).
		Model(&model.AdminRule{}).
		Where("id IN ? AND pid IS NOT NULL", ids).
		Distinct().
		Pluck("pid", &pids).Error
	if err != nil {
		return nil, err
	}
	return pids, nil
}

// TitleMapByIDs 返回 ids 中指定 types 的规则的 id -> title 映射
func (r *AdminRuleRepository) TitleMapByIDs(c *gin.Context, ids []uint, types ...string) (map[uint]string, error) {
	if len(ids) == 0 {
		return map[uint]string{}, nil
	}
	q := gorm.G[model.AdminRule](r.DB()).Where("id IN ?", ids)
	if len(types) > 0 {
		q = q.Where("type IN ?", types)
	}
	rules, err := q.Find(c.Request.Context())
	if err != nil {
		return nil, err
	}
	result := make(map[uint]string, len(rules))
	for _, rule := range rules {
		result[rule.ID] = rule.Title
	}
	return result, nil
}
