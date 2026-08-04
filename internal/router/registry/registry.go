package registry

import "github.com/gin-gonic/gin"

// Routes 已注册的路由模块列表，由各子模块通过 init() 填充
var Routes []func(*gin.Engine)

// Register 子模块通过 init() 调用，接受一个 gin.Engine，将路由注册函数追加到 Routes
func Register(fn func(*gin.Engine)) {
	Routes = append(Routes, fn)
}

// AdminRoutes 已注册的后台路由模块列表，由各子模块通过 init() 填充
var AdminRoutes []func(*gin.RouterGroup)

// RegisterAdmin 子模块通过 init() 调用，接受一个 gin.RouterGroup，将后台路由注册函数追加到 AdminRoutes
// admin 单独一个路由注册函数而不是使用 Register，以便实现自定义后台入口，同时方便挂载后台级中间件
func RegisterAdmin(fn func(*gin.RouterGroup)) {
	AdminRoutes = append(AdminRoutes, fn)
}
