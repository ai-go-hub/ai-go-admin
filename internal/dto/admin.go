package dto

import "github.com/ai-go-hub/ai-go-admin/internal/model"

// AdminSession 管理员会话信息，由认证中间件注入到请求上下文
type AdminSession struct {
	*model.Admin
	Token string `json:"token"`
}
