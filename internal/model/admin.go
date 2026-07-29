package model

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID           uint            `gorm:"comment:ID;primarykey;autoIncrement" json:"id"`
	Username     string          `gorm:"comment:用户名;type:varchar(64);not null;default:''" json:"username" binding:"required"`
	Nickname     string          `gorm:"comment:昵称;type:varchar(64);not null;default:''" json:"nickname" binding:"required"`
	Avatar       *string         `gorm:"comment:头像;type:varchar(255)" json:"avatar"`
	Email        *string         `gorm:"comment:邮箱;type:varchar(128)" json:"email"`
	Mobile       *string         `gorm:"comment:手机号;type:varchar(16)" json:"mobile"`
	LoginFailure uint            `gorm:"comment:连续登录失败次数;not null;default:0" json:"-"`
	LastLoginAt  *time.Time      `gorm:"comment:上次登录时间" json:"last_login_at"`
	LastLoginIP  *string         `gorm:"comment:上次登录IP;type:varchar(64)" json:"last_login_ip"`
	Password     string          `gorm:"comment:密码;type:varchar(255);not null;default:''" json:"-"`
	Bio          *string         `gorm:"comment:个人简介;type:varchar(255)" json:"bio"`
	Status       string          `gorm:"comment:状态:enable=启用,disable=禁用;type:varchar(64);not null;default:'enable'" json:"status" binding:"required"`
	UpdatedAt    *time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	CreatedAt    *time.Time      `gorm:"comment:创建时间" json:"created_at"`
	DeletedAt    *gorm.DeletedAt `gorm:"comment:删除时间;index" json:"-"`
}

func (Admin) TableName() string {
	return "admins"
}

// AdminRule 菜单和权限规则模型
type AdminRule struct {
	ID        uint       `gorm:"comment:ID;primarykey;autoIncrement" json:"id"`
	Pid       *uint      `gorm:"comment:上级规则" json:"pid"`
	Type      string     `gorm:"comment:规则类型:dir=规则目录,menu=菜单项,node=权限节点;type:varchar(64);not null;default:''" binding:"required,oneof=dir menu node" json:"type"`
	Title     string     `gorm:"comment:规则标题;type:varchar(64);not null;default:''" binding:"required" json:"title"`
	Name      string     `gorm:"comment:规则名称;type:varchar(64);not null;default:''" binding:"required" json:"name"`
	Path      *string    `gorm:"comment:菜单路由路径;type:varchar(255);default:''" json:"path"`
	Icon      *string    `gorm:"comment:菜单图标;type:varchar(64);default:''" json:"icon"`
	OpenType  *string    `gorm:"comment:菜单打开方式:tab=选项卡,link=链接,iframe=Iframe;type:varchar(64);default:''" json:"open_type"`
	URL       *string    `gorm:"comment:菜单URL;type:varchar(255);default:''" json:"url"`
	Component *string    `gorm:"comment:菜单组件路径;type:varchar(255);default:''" json:"component"`
	Keepalive *uint8     `gorm:"comment:缓存:0=关闭,1=开启;default:0" json:"keepalive"`
	Extend    *string    `gorm:"comment:扩展属性:add_route_only=只添加为路由,add_menu_only=只添加为菜单;type:varchar(64);default:''" json:"extend"`
	Remark    *string    `gorm:"comment:备注;type:varchar(255);default:''" json:"remark"`
	Weigh     *int       `gorm:"comment:权重;default:0" json:"weigh"`
	Status    *uint8     `gorm:"comment:状态:0=禁用,1=启用;not null;default:1" binding:"required" json:"status"`
	UpdatedAt *time.Time `gorm:"comment:更新时间" json:"updated_at"`
	CreatedAt *time.Time `gorm:"comment:创建时间" json:"created_at"`
}

func (AdminRule) TableName() string {
	return "admin_rules"
}

// AfterCreate 新增后，若 Weigh 未指定则回填为主键 ID
func (r *AdminRule) AfterCreate(tx *gorm.DB) error {
	if r.Weigh != nil && *r.Weigh != 0 {
		return nil
	}
	weigh := int(r.ID)
	r.Weigh = &weigh
	return tx.Model(r).Update("weigh", r.ID).Error
}

// ToMap 将 AdminRule 转换为 map
// 菜单规则使用较频繁，使用硬编码转换性能最高
func (r *AdminRule) ToMap() map[string]any {
	pid := uint(0)
	if r.Pid != nil {
		pid = *r.Pid
	}
	return map[string]any{
		"id":         r.ID,
		"pid":        pid,
		"type":       r.Type,
		"title":      r.Title,
		"name":       r.Name,
		"path":       r.Path,
		"icon":       r.Icon,
		"open_type":  r.OpenType,
		"url":        r.URL,
		"component":  r.Component,
		"keepalive":  r.Keepalive,
		"extend":     r.Extend,
		"remark":     r.Remark,
		"weigh":      r.Weigh,
		"status":     r.Status,
		"updated_at": r.UpdatedAt,
		"created_at": r.CreatedAt,
	}
}

// AdminGroup 管理员角色组模型
type AdminGroup struct {
	ID        uint       `gorm:"comment:ID;primarykey;autoIncrement" json:"id"`
	Pid       *uint      `gorm:"comment:上级分组;default:0" json:"pid"`
	Name      string     `gorm:"comment:组名;type:varchar(64);not null;default:''" binding:"required" json:"name"`
	Rules     *string    `gorm:"comment:权限规则ID集;type:text" json:"rules"`
	Status    *uint8     `gorm:"comment:状态:0=禁用,1=启用;not null;default:1" binding:"required" json:"status"`
	UpdatedAt *time.Time `gorm:"comment:更新时间" json:"updated_at"`
	CreatedAt *time.Time `gorm:"comment:创建时间" json:"created_at"`
}

func (AdminGroup) TableName() string {
	return "admin_groups"
}

// ToMap 将 AdminGroup 转换为 map
// 使用硬编码转换性能最高
func (g *AdminGroup) ToMap() map[string]any {
	pid := uint(0)
	if g.Pid != nil {
		pid = *g.Pid
	}
	return map[string]any{
		"id":         g.ID,
		"pid":        pid,
		"name":       g.Name,
		"rules":      g.Rules,
		"status":     g.Status,
		"updated_at": g.UpdatedAt,
		"created_at": g.CreatedAt,
	}
}

// AdminGroupAccess 管理员分组映射模型
type AdminGroupAccess struct {
	UID     uint `gorm:"comment:管理员ID;not null;index" json:"uid"`
	GroupID uint `gorm:"comment:分组ID;not null;index" json:"group_id"`
}

func (AdminGroupAccess) TableName() string {
	return "admin_group_access"
}
