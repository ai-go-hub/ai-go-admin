package model

import "time"

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
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	ExpiredAt time.Time `gorm:"comment:过期时间;not null;index" json:"expired_at"`
}

// TableName 指定 Captcha 模型表名
func (Captcha) TableName() string {
	return "captchas"
}

// Config 系统配置模型
type Config struct {
	ID          uint   `gorm:"comment:ID;primarykey;autoIncrement" json:"id"`
	Name        string `gorm:"comment:变量名;type:varchar(50);uniqueIndex;not null;default:''" json:"name"`
	Group       string `gorm:"comment:分组;type:varchar(50);not null;default:''" json:"group"`
	Title       string `gorm:"comment:变量标题;type:varchar(50);not null;default:''" json:"title"`
	Tip         string `gorm:"comment:变量描述;type:varchar(100);not null;default:''" json:"tip"`
	Type        string `gorm:"comment:变量输入组件类型;type:varchar(50);not null;default:''" json:"type"`
	Value       string `gorm:"comment:变量值;type:text" json:"value"`
	Content     string `gorm:"comment:字典数据;type:text" json:"content"`
	Rule        string `gorm:"comment:验证规则;type:varchar(100);not null;default:''" json:"rule"`
	Extend      string `gorm:"comment:扩展属性;type:varchar(255);not null;default:''" json:"extend"`
	InputExtend string `gorm:"comment:输入框扩展属性;type:varchar(255);not null;default:''" json:"input_extend"`
	AllowDel    uint8  `gorm:"comment:允许删除:0=否,1=是;not null;default:0" json:"allow_del"`
	Weigh       int    `gorm:"comment:权重;not null;default:0" json:"weigh"`
}

// TableName 指定 Config 模型表名
func (Config) TableName() string {
	return "configs"
}
