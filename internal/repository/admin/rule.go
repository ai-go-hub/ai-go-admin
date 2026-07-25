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
