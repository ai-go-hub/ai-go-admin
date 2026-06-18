package model

import "time"

func init() {
	Register(&Token{}, &Captcha{})
}

// Token 令牌模型，用于存储各类用户令牌
type Token struct {
	Token     string    `gorm:"comment:令牌;type:varchar(64);primaryKey" json:"-"`
	Type      string    `gorm:"comment:令牌类型;type:varchar(32);not null" json:"type"`
	UserID    uint      `gorm:"comment:用户ID;not null;index" json:"user_id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	ExpiredAt time.Time `gorm:"comment:过期时间;not null;index" json:"expired_at"`
}

// TableName 指定 Token 模型表名
func (Token) TableName() string {
	return "tokens"
}

// Captcha 验证码模型
type Captcha struct {
	Key       string    `gorm:"comment:验证码查询键;type:varchar(64);primaryKey" json:"key"`
	Code      string    `gorm:"comment:验证码值（加密后）;type:varchar(255)" json:"-"`
	Info      string    `gorm:"comment:验证码详细信息;type:text" json:"-"`
	ExpiredAt time.Time `gorm:"comment:过期时间;not null;index" json:"expired_at"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
}

// TableName 指定 Captcha 模型表名
func (Captcha) TableName() string {
	return "captchas"
}
