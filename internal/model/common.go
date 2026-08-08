package model

import (
	"time"

	"gorm.io/datatypes"
)

// Token 令牌模型
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
	ID          uint    `gorm:"comment:ID;primarykey;autoIncrement" json:"id"`
	Name        string  `gorm:"comment:变量名;type:varchar(64);uniqueIndex;not null;default:''" json:"name"`
	Group       string  `gorm:"comment:分组;type:varchar(64);not null;default:''" json:"group"`
	Title       string  `gorm:"comment:变量标题;type:varchar(64);not null;default:''" json:"title"`
	Type        string  `gorm:"comment:变量输入组件类型;type:varchar(64);not null;default:''" json:"type"`
	Tip         *string `gorm:"comment:变量输入提示;type:varchar(128)" json:"tip"`
	Value       *string `gorm:"comment:变量值;type:text" json:"value"`
	Dict        *string `gorm:"comment:字典数据;type:text" json:"dict"`
	Rule        *string `gorm:"comment:验证规则;type:varchar(128)" json:"rule"`
	InputExtend *string `gorm:"comment:输入框扩展属性;type:varchar(255)" json:"input_extend"`
	AllowDel    uint8   `gorm:"comment:允许删除:0=否,1=是;not null;default:1" json:"allow_del"`
	Weigh       int     `gorm:"comment:权重;not null;default:0" json:"weigh"`
}

// TableName 指定 Config 模型表名
func (Config) TableName() string {
	return "configs"
}

// Attachment 附件模型
type Attachment struct {
	ID           uint      `gorm:"comment:ID;primarykey;autoIncrement" json:"id"`
	Topic        string    `gorm:"comment:主题分类;type:varchar(64);not null;default:''" json:"topic"`
	UserID       uint      `gorm:"comment:上传用户ID;not null;default:0" json:"user_id"`
	UserType     string    `gorm:"comment:上传用户身份类型;type:varchar(64);not null;default:''" json:"user_type"`
	URL          string    `gorm:"comment:存储路径;type:varchar(255);not null;default:''" json:"url"`
	Name         string    `gorm:"comment:原始名称;type:varchar(255);not null;default:''" json:"name"`
	Size         int64     `gorm:"comment:大小;not null;default:0" json:"size"`
	Mimetype     string    `gorm:"comment:MIME类型;type:varchar(64);not null;default:''" json:"mimetype"`
	Quote        int64     `gorm:"comment:上传（引用）次数;not null;default:0" json:"quote"`
	Driver       string    `gorm:"comment:存储驱动;type:varchar(64);not null;default:''" json:"driver"`
	Sha1         string    `gorm:"comment:SHA1编码;type:varchar(64);not null;default:'';uniqueIndex" json:"sha1"`
	CreatedAt    time.Time `gorm:"comment:创建时间" json:"created_at"`
	LastUploadAt time.Time `gorm:"comment:最后上传时间" json:"last_upload_at"`
}

// TableName 指定 Attachment 模型表名
func (Attachment) TableName() string {
	return "attachments"
}

// Area 省份地区模型
type Area struct {
	ID    uint   `gorm:"comment:ID;primarykey;autoIncrement" json:"id"`
	Pid   uint   `gorm:"comment:上级ID;not null;default:0;index" json:"pid"`
	Name  string `gorm:"comment:名称;type:varchar(64);not null;default:''" json:"name"`
	Level uint8  `gorm:"comment:等级;not null;default:0;index" json:"level"`
	Code  string `gorm:"comment:行政区划代码;type:varchar(32);not null;default:''" json:"code"`
	Zip   string `gorm:"comment:邮编;type:varchar(32);not null;default:''" json:"zip"`
}

// TableName 指定 Area 模型表名
func (Area) TableName() string {
	return "areas"
}

// CrudLog CRUD 记录模型
type CrudLog struct {
	ID        uint           `gorm:"comment:ID;primarykey;autoIncrement" json:"id"`
	Name      string         `gorm:"comment:表名;type:varchar(255);not null;default:''" json:"name"`
	Comment   string         `gorm:"comment:注释;type:varchar(255);not null;default:''" json:"comment"`
	Table     datatypes.JSON `gorm:"comment:数据表数据;type:jsonb" json:"table"`
	Fields    datatypes.JSON `gorm:"comment:字段数据;type:jsonb" json:"fields"`
	Status    string         `gorm:"comment:状态:deleted=已删除,succeeded=成功,failed=失败,generating=生成中;type:varchar(64);not null;default:''" json:"status"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
}

// TableName 指定 CrudLog 模型表名
func (CrudLog) TableName() string {
	return "crud_logs"
}
