package router

import (
	// 空白导入触发子目录 init() 自动注册路由
	_ "github.com/ai-go-hub/ai-go-admin/internal/router/admin"
	_ "github.com/ai-go-hub/ai-go-admin/internal/router/common"

	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"

	"github.com/gin-gonic/gin"
)

// Setup 遍历所有已注册的路由模块，传入 Engine 完成注册
func Setup(r *gin.Engine) {
	for _, fn := range registry.Routes {
		fn(r)
	}
}
