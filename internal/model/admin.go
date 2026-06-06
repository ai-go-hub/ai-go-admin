package model

import (
	"time"

	"gorm.io/gorm"
)

func init() {
	Register(&Admin{})
}

type Admin struct {
	ID           uint           `gorm:"comment:ID;primarykey;autoIncrement"`
	Username     string         `gorm:"comment:用户名;type:varchar(64)"`
	Nickname     string         `gorm:"comment:昵称;type:varchar(64)"`
	Avatar       string         `gorm:"comment:头像;type:varchar(255)"`
	Email        string         `gorm:"comment:邮箱;type:varchar(128)"`
	Mobile       string         `gorm:"comment:手机号;type:varchar(16)"`
	LoginFailure uint           `gorm:"comment:连续登录失败次数"`
	LastLoginAt  time.Time      `gorm:"comment:上次登录时间"`
	LastLoginIp  string         `gorm:"comment:上次登录IP;type:varchar(64)"`
	Password     string         `gorm:"comment:密码;type:varchar(255)"`
	Bio          string         `gorm:"comment:个人简介;type:varchar(255)"`
	Status       string         `gorm:"comment:状态:enable=启用,disable=禁用;type:varchar(64)"`
	UpdatedAt    time.Time      `gorm:"comment:更新时间"`
	CreatedAt    time.Time      `gorm:"comment:创建时间"`
	DeletedAt    gorm.DeletedAt `gorm:"comment:删除时间;index"`
}
