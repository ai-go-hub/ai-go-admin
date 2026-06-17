package model

import "time"

func init() {
	Register(&Token{})
}

// Token 令牌模型，用于存储各类用户令牌
type Token struct {
	Token     string    `gorm:"comment:令牌;type:varchar(64);primaryKey" json:"-"`
	Type      string    `gorm:"comment:令牌类型;type:varchar(32);not null" json:"type"`
	UserID    uint      `gorm:"comment:用户ID;not null;index" json:"user_id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	ExpiredAt time.Time `gorm:"comment:过期时间;not null;index" json:"expired_at"`
}

// TableName 指定表名
func (Token) TableName() string {
	return "tokens"
}
