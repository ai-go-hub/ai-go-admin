package common

import (
	handlerCommon "github.com/ai-go-hub/ai-go-admin/internal/handler/common"
	repoCommon "github.com/ai-go-hub/ai-go-admin/internal/repository/common"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"
	svcCommon "github.com/ai-go-hub/ai-go-admin/internal/service/common"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.Register(func(r *gin.Engine) {
		repo := repoCommon.NewCaptchaRepository()
		svc := svcCommon.NewCaptchaService(repo)
		h := handlerCommon.NewCaptchaHandler(svc)

		group := r.Group("/common/captcha")
		h.RegisterRoutes(group)
	})
}
