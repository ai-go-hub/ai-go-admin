package crud

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/handler"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/service"
	svcCrud "github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"

	"github.com/gin-gonic/gin"
)

// CrudLogListItem CRUD 记录列表项，在模型基础上附加 label 展示字段
type CrudLogListItem struct {
	model.CrudLog
	Label           string                        `json:"label"`
	RouterBasicData dto.GenerateFileBasicDataInfo `json:"router_basic_data"`
}

// CrudLogHandler CRUD 记录控制器，嵌入通用控制器
type CrudLogHandler struct {
	*handler.Handler[model.CrudLog]
	svc *svcCrud.CrudLogService
}

// NewCrudLogHandler 创建 CRUD 记录控制器实例
func NewCrudLogHandler(svc *svcCrud.CrudLogService) *CrudLogHandler {
	return &CrudLogHandler{
		Handler: handler.NewHandler(svc,
			handler.WithAdapter(handler.Adapter{
				// 列表额外返回 label 字段: name - comment
				List: func(ctx context.Context, list any, opts service.Options) (any, error) {
					items, ok := list.([]model.CrudLog)
					if !ok {
						return nil, errors.New("列表数据类型错误")
					}

					result := make([]CrudLogListItem, 0, len(items))
					for _, item := range items {
						var tableData dto.CRUDTable
						_ = json.Unmarshal(item.Table, &tableData)
						result = append(result, CrudLogListItem{
							CrudLog:         item,
							Label:           item.Name + " - " + item.Comment,
							RouterBasicData: svcCrud.ParseGenerateFileBasicData("router", tableData.RouterFile),
						})
					}
					return result, nil
				},
			}),
			handler.WithOmitFields(handler.ActionFields{
				Update: []string{"id"},
			}),
		),
		svc: svc,
	}
}

// RegisterRoutes 注册路由
func (h *CrudLogHandler) RegisterRoutes(group *gin.RouterGroup) {
	handler.RegisterBaseRoutes(h, group)
}
