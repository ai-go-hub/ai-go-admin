package permission

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"
	"github.com/ai-go-hub/ai-go-admin/internal/model"

	"gorm.io/gorm"
)

// Permission 权限规则管理器，通过管理员分组和权限规则表进行权限判断
type Permission struct{}

// New 创建 Permission 实例
func New() *Permission {
	return &Permission{}
}

// GetGroups 根据管理员 ID 获取管理员的分组列表（仅返回启用状态的分组）
func (p *Permission) GetGroups(ctx context.Context, adminId uint) ([]model.AdminGroup, error) {
	// 1. 查询管理员的分组关系
	accesses, err := gorm.G[model.AdminGroupAccess](database.DB()).Where("uid = ?", adminId).Find(ctx)
	if err != nil {
		return nil, err
	}

	if len(accesses) == 0 {
		return nil, nil
	}

	// 2. 收集分组 ID
	groupIDs := make([]uint, len(accesses))
	for i, acc := range accesses {
		groupIDs[i] = acc.GroupID
	}

	// 3. 查询启用状态的分组
	groups, err := gorm.G[model.AdminGroup](database.DB()).Where("id IN ? AND status = 1", groupIDs).Find(ctx)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// GetRuleIds 根据管理员 ID 获取管理员拥有权限的全部规则 ID（去重）
func (p *Permission) GetRuleIds(ctx context.Context, adminId uint) ([]uint, error) {
	groups, err := p.GetGroups(ctx, adminId)
	if err != nil {
		return nil, err
	}

	return p.ruleIdsFromGroups(ctx, groups)
}

// GetRules 根据管理员 ID 获取管理员拥有的权限规则列表（仅返回启用状态的规则）
func (p *Permission) GetRules(ctx context.Context, adminId uint) ([]model.AdminRule, error) {
	groups, err := p.GetGroups(ctx, adminId)
	if err != nil {
		return nil, err
	}

	if slices.ContainsFunc(groups, isSuperAdminGroup) {
		return gorm.G[model.AdminRule](database.DB()).Where("status = 1").Find(ctx)
	}

	ruleIDs, err := p.ruleIdsFromGroups(ctx, groups)
	if err != nil {
		return nil, err
	}

	if len(ruleIDs) == 0 {
		return nil, nil
	}

	rules, err := gorm.G[model.AdminRule](database.DB()).Where("id IN ? AND status = 1", ruleIDs).Find(ctx)
	if err != nil {
		return nil, err
	}

	return rules, nil
}

// Check 检查管理员是否拥有指定名称的权限规则
func (p *Permission) Check(ctx context.Context, adminId uint, ruleName string) (bool, error) {
	groups, err := p.GetGroups(ctx, adminId)
	if err != nil {
		return false, err
	}

	if slices.ContainsFunc(groups, isSuperAdminGroup) {
		return true, nil
	}

	ruleIDs, err := p.ruleIdsFromGroups(ctx, groups)
	if err != nil {
		return false, err
	}

	if len(ruleIDs) == 0 {
		return false, nil
	}

	_, err = gorm.G[model.AdminRule](database.DB()).Where("id IN ? AND name = ? AND status = 1", ruleIDs, ruleName).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// isSuperAdminGroup 判断分组是否拥有全部权限（Rules 字段值为 "*"）
func isSuperAdminGroup(group model.AdminGroup) bool {
	return group.Rules != nil && *group.Rules == "*"
}

// ruleIdsFromGroups 从分组列表中提取规则 ID，遇到通配符分组则返回全部启用规则的 ID
func (p *Permission) ruleIdsFromGroups(ctx context.Context, groups []model.AdminGroup) ([]uint, error) {
	if slices.ContainsFunc(groups, isSuperAdminGroup) {
		return p.allRuleIDs(ctx)
	}

	ruleIDSet := make(map[uint]struct{})
	for _, group := range groups {
		if group.Rules == nil || *group.Rules == "" {
			continue
		}

		for part := range strings.SplitSeq(*group.Rules, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			id, err := strconv.ParseUint(part, 10, 32)
			if err != nil {
				continue
			}
			ruleIDSet[uint(id)] = struct{}{}
		}
	}

	ruleIDs := make([]uint, 0, len(ruleIDSet))
	for id := range ruleIDSet {
		ruleIDs = append(ruleIDs, id)
	}

	return ruleIDs, nil
}

// allRuleIDs 返回全部启用规则的 ID
func (p *Permission) allRuleIDs(ctx context.Context) ([]uint, error) {
	rules, err := gorm.G[model.AdminRule](database.DB()).Where("status = 1").Find(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, len(rules))
	for i, rule := range rules {
		ids[i] = rule.ID
	}

	return ids, nil
}
