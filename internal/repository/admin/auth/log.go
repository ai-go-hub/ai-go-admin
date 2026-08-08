package auth

import (
	"context"
	"fmt"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"
)

// AdminLogRepository 管理员日志仓储，嵌入通用仓储
type AdminLogRepository struct {
	*repository.Repository[model.AdminLog]
}

// NewAdminLogRepository 创建管理员日志仓储实例
func NewAdminLogRepository() *AdminLogRepository {
	return &AdminLogRepository{
		Repository: repository.NewRepository[model.AdminLog](),
	}
}

// GetRuleTitle 根据规则 name 获取格式化标题
// 使用自连接（LEFT JOIN）在一条 SQL 中同时取出当前规则标题及其父级菜单标题
// 若存在父级，返回 "父级标题 - 规则标题" 格式，否则返回规则标题
// 若规则不存在，返回 name 本身作为兜底
func (r *AdminLogRepository) GetRuleTitle(ctx context.Context, name string) string {
	db := r.DB().WithContext(ctx)

	// 实例化 AdminRule 仓储，复用基类 Schema() 获取完整表名
	sch, _ := repository.NewRepository[model.AdminRule]().Schema()

	var result struct {
		Title       string
		ParentTitle string
	}
	err := db.Raw(
		fmt.Sprintf(
			"SELECT c.title, COALESCE(p.title, '') AS parent_title FROM %s AS c LEFT JOIN %s AS p ON p.id = c.pid WHERE c.name = ?",
			sch.Table, sch.Table,
		),
		name,
	).Scan(&result).Error

	if err != nil || result.Title == "" {
		return name
	}
	if result.ParentTitle != "" {
		return result.ParentTitle + " - " + result.Title
	}
	return result.Title
}
