package model

import "time"

func init() {
	Register(&User{})
}

type User struct {
	ID        uint
	Name      string `gorm:"type:varchar(64);comment:用户名"`
	Email     string `gorm:"type:varchar(128);comment:邮箱"`
	Mobile    string `gorm:"type:varchar(16);comment:手机号"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
