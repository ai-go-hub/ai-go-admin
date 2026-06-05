package registry

import "github.com/gin-gonic/gin"

// Routes 已注册的路由模块列表，由各子模块通过 init() 填充
var Routes []func(*gin.Engine)

// Register 子模块通过 init() 调用，将路由注册函数追加到 Routes
func Register(fn func(*gin.Engine)) {
	Routes = append(Routes, fn)
}
