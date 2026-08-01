package routine

import (
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	repoCommon "github.com/ai-go-hub/ai-go-admin/internal/repository/common"
	"github.com/ai-go-hub/ai-go-admin/internal/service"
)

// AttachmentService 附件服务
type AttachmentService struct {
	service.IService[model.Attachment]
	repo *repoCommon.AttachmentRepository
}

// NewAttachmentService 创建附件服务实例
func NewAttachmentService(repo *repoCommon.AttachmentRepository) *AttachmentService {
	return &AttachmentService{
		IService: service.NewService(repo),
		repo:     repo,
	}
}
