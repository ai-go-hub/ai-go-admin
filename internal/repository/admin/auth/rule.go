package auth

import (
	"context"
	"errors"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"

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

// BatchCreate 批量创建菜单规则树，同名规则已存在则跳过，支持递归创建子级规则
func (r *AdminRuleRepository) BatchCreate(ctx context.Context, menus []*model.AdminRule, parentID *uint) error {
	for _, menu := range menus {
		if menu == nil {
			continue
		}
		menu.Pid = parentID

		// 同名规则已存在则跳过创建，其子级继续以该规则为父递归创建
		existing, err := gorm.G[model.AdminRule](r.DB()).Where("name = ?", menu.Name).First(ctx)
		if err == nil {
			if len(menu.Children) > 0 {
				childParentID := existing.ID
				if err := r.BatchCreate(ctx, menu.Children, &childParentID); err != nil {
					return err
				}
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 创建当前规则后，以其 ID 为父递归创建子级
		if err := r.Create(ctx, menu, repository.Options{}); err != nil {
			return err
		}
		if len(menu.Children) > 0 {
			childParentID := menu.ID
			if err := r.BatchCreate(ctx, menu.Children, &childParentID); err != nil {
				return err
			}
		}
	}
	return nil
}

// BatchDelete 批量删除菜单规则
//
// cleanup=true 时进行干净清理: 递归删除 ids 指定菜单的所有子级、删除 ids 指定的菜单，
// 并逐级递归删除因删除而变空的上级目录（Type=dir）
func (r *AdminRuleRepository) BatchDelete(ctx context.Context, ids []uint, cleanup bool) error {
	if len(ids) == 0 {
		return nil
	}

	// 收集待删除规则 id（ids + 递归子级），并记录其直接父级（供清理上级空目录）
	delIDs := make([]uint, 0, len(ids))
	seen := make(map[uint]struct{}, len(ids))
	var parentIDs []uint
	stack := append([]uint{}, ids...)
	for len(stack) > 0 {
		id := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		delIDs = append(delIDs, id)

		childIDs, err := r.ChildIDsByPids(ctx, []uint{id})
		if err != nil {
			return err
		}
		stack = append(stack, childIDs...)
	}

	// 删除前先收集直接父级（删除后其行已不存在）
	if cleanup {
		var err error
		parentIDs, err = r.DistinctPidsByIDs(ctx, delIDs)
		if err != nil {
			return err
		}
	}

	if err := r.Delete(ctx, repository.Options{PrimaryKeyValues: util.UintsToStrs(delIDs)}); err != nil {
		return err
	}
	if !cleanup {
		return nil
	}

	// 逐级向上删除已无子级的目录
	for _, pid := range parentIDs {
		if err := r.deleteEmptyDirChain(ctx, pid); err != nil {
			return err
		}
	}
	return nil
}

// deleteEmptyDirChain 逐级向上删除无子级的菜单规则
func (r *AdminRuleRepository) deleteEmptyDirChain(ctx context.Context, id uint) error {
	for id > 0 {
		rule, err := gorm.G[model.AdminRule](r.DB()).Where("id = ?", id).First(ctx)
		if err != nil {
			return nil
		}
		if rule.Type != "dir" {
			return nil
		}
		childIDs, err := r.ChildIDsByPids(ctx, []uint{id})
		if err != nil {
			return err
		}
		if len(childIDs) > 0 {
			return nil
		}
		if err := r.Delete(ctx, repository.Options{PrimaryKeyValues: util.UintsToStrs([]uint{id})}); err != nil {
			return err
		}
		if rule.Pid == nil {
			return nil
		}
		id = *rule.Pid
	}
	return nil
}
