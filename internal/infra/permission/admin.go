package permission

import (
	"context"
	"slices"
	"strconv"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"

	"gorm.io/gorm"
)

// 规则状态常量，与 model.AdminRule.Status 对齐
const (
	RuleStatusDisabled uint8 = 0 // 禁用
	RuleStatusEnabled  uint8 = 1 // 启用
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

// IsSuperAdmin 是否超级管理员
func (p *Permission) IsSuperAdmin(ctx context.Context, adminId uint) (bool, error) {
	groups, err := p.GetGroups(ctx, adminId)
	if err != nil {
		return false, err
	}
	return slices.ContainsFunc(groups, isSuperAdminGroup), nil
}

// GetRuleIds 根据管理员用户 ID 获取管理员的权限规则 ID（去重）
func (p *Permission) GetRuleIds(ctx context.Context, adminId uint, ruleStatus *uint8) ([]uint, error) {
	groups, err := p.GetGroups(ctx, adminId)
	if err != nil {
		return nil, err
	}

	return p.ruleIdsFromGroups(ctx, groups, ruleStatus)
}

// GetRules 根据管理员 ID 获取管理员拥有的权限规则列表
func (p *Permission) GetRules(ctx context.Context, adminId uint, ruleStatus *uint8) ([]model.AdminRule, error) {
	groups, err := p.GetGroups(ctx, adminId)
	if err != nil {
		return nil, err
	}

	ruleIDs, err := p.ruleIdsFromGroups(ctx, groups, ruleStatus)
	if err != nil {
		return nil, err
	}

	if len(ruleIDs) == 0 {
		return nil, nil
	}

	q := gorm.G[model.AdminRule](database.DB()).Where("id IN ?", ruleIDs)
	if ruleStatus != nil {
		q = q.Where("status = ?", *ruleStatus)
	}
	rules, err := q.Order("weigh desc").Find(ctx)
	if err != nil {
		return nil, err
	}

	return rules, nil
}

// Check 检查管理员是否拥有指定名称的权限规则
// ruleNames 为待检查的权限节点名称列表，op 取值 "OR" 或 "AND"（默认），
// OR 表示任意一个命中即返回 true，AND 表示必须全部命中才返回 true
func (p *Permission) Check(ctx context.Context, adminId uint, ruleNames []string, op string) (bool, error) {
	if len(ruleNames) == 0 {
		return false, nil
	}

	// 将 op 转小写后，只能是 or 或 and
	op = strings.ToLower(op)
	if op != "or" {
		op = "and"
	}

	groups, err := p.GetGroups(ctx, adminId)
	if err != nil {
		return false, err
	}

	// 超管则提前短路
	if slices.ContainsFunc(groups, isSuperAdminGroup) {
		return true, nil
	}

	ruleIDs, err := p.ruleIdsFromGroups(ctx, groups, util.ToPtr(RuleStatusEnabled))
	if err != nil {
		return false, err
	}

	if len(ruleIDs) == 0 {
		return false, nil
	}

	rules, err := gorm.G[model.AdminRule](database.DB()).
		Where("id IN ? AND name IN ? AND status = ?", ruleIDs, ruleNames, RuleStatusEnabled).
		Find(ctx)
	if err != nil {
		return false, err
	}

	if op == "or" {
		return len(rules) > 0, nil
	}
	return len(rules) == len(ruleNames), nil
}

// isSuperAdminGroup 判断分组是否拥有全部权限（即超管，Rules 字段值为 "*"）
func isSuperAdminGroup(group model.AdminGroup) bool {
	return group.Rules != nil && *group.Rules == "*"
}

// ruleIdsFromGroups 从管理分组列表中提取权限规则 ID，
// 遇到通配符分组（超管）则返回全部规则 ID
func (p *Permission) ruleIdsFromGroups(ctx context.Context, groups []model.AdminGroup, ruleStatus *uint8) ([]uint, error) {
	if slices.ContainsFunc(groups, isSuperAdminGroup) {
		return p.AllRuleIDs(ctx, ruleStatus)
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

// AllRuleIDs 返回全部规则 ID，
// ruleStatus 为 nil 时不按状态过滤
func (p *Permission) AllRuleIDs(ctx context.Context, ruleStatus *uint8) ([]uint, error) {
	var scopes []func(*gorm.Statement)
	if ruleStatus != nil {
		s := *ruleStatus // 先捕获值
		scopes = append(scopes, func(stmt *gorm.Statement) {
			stmt.Where("status = ?", s)
		})
	}

	rules, err := gorm.G[model.AdminRule](database.DB()).Scopes(scopes...).Select("id").Find(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, len(rules))
	for i, rule := range rules {
		ids[i] = rule.ID
	}

	return ids, nil
}

// BuildCheckPath 从 fullPath 获取权限节点名称
// 去掉 "/admin/" 前缀，去掉 params（以 ":" 开头的路径段），然后以 "/" 连接剩余部分
func BuildCheckPath(fullPath string) string {
	path := strings.TrimPrefix(fullPath, "/admin/")
	parts := strings.Split(path, "/")
	filtered := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.HasPrefix(p, ":") {
			break
		}
		if p != "" {
			filtered = append(filtered, p)
		}
	}
	return strings.Join(filtered, "/")
}
