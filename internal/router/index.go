package router

import (
	"net/http"

	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.Register(func(r *gin.Engine) {
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "index pong"})
		})
	})
}
