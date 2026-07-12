package router

import (
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.Register(func(r *gin.Engine) {
		// 静态资源托管
		r.Static("/static", "./static")
	})
}
