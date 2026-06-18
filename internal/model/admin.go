package model

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	Register(&Admin{})
}

type Admin struct {
	ID           uint           `gorm:"comment:ID;primarykey;autoIncrement" json:"id"`
	Username     string         `gorm:"comment:用户名;type:varchar(64)" json:"username"`
	Nickname     string         `gorm:"comment:昵称;type:varchar(64)" json:"nickname"`
	Avatar       string         `gorm:"comment:头像;type:varchar(255)" json:"avatar"`
	Email        string         `gorm:"comment:邮箱;type:varchar(128)" json:"email"`
	Mobile       string         `gorm:"comment:手机号;type:varchar(16)" json:"mobile"`
	LoginFailure uint           `gorm:"comment:连续登录失败次数;not null;default:0" json:"-"`
	LastLoginAt  time.Time      `gorm:"comment:上次登录时间" json:"last_login_at"`
	LastLoginIP  string         `gorm:"comment:上次登录IP;type:varchar(64)" json:"last_login_ip"`
	Password     string         `gorm:"comment:密码;type:varchar(255)" json:"-"`
	Bio          string         `gorm:"comment:个人简介;type:varchar(255)" json:"bio"`
	Status       string         `gorm:"comment:状态:enable=启用,disable=禁用;type:varchar(64)" json:"status"`
	UpdatedAt    time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	CreatedAt    time.Time      `gorm:"comment:创建时间" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"comment:删除时间;index" json:"-"`
}

func (Admin) TableName() string {
	return "admins"
}
