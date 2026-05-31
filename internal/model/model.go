package model

var models []any

// 注册需要自动迁移的模型，各模型文件内可通过 init() 机制实现自动调用
func Register(m any) {
	models = append(models, m)
}

// 返回所有已注册的模型，供 database.Init 的 AutoMigrate 使用
func All() []any {
	return models
}
