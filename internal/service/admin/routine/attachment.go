package routine

import (
	"context"
	"fmt"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/upload"
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

// Delete 覆写通用删除方法: 在删除数据库记录前先通过驱动删除物理文件
func (s *AttachmentService) Delete(ctx context.Context, opts service.Options) error {
	// 查询待删除的附件记录，获取 Driver 和 URL
	pkField, err := s.repo.PrimaryKeyField()
	if err != nil {
		return err
	}

	wheres := []service.WhereGroup{{
		Wheres: []service.Where{{
			Field:    pkField,
			Operator: "IN",
			Value:    opts.PrimaryKeyValues,
		}},
	}}

	attachments, err := s.IService.List(ctx, service.Options{Wheres: wheres})
	if err != nil {
		return err
	}

	// 遍历删除物理文件
	for _, att := range attachments {
		d, err := upload.NewDriver(att.Driver)
		if err != nil {
			return fmt.Errorf("驱动 %s 错误，文件删除失败: %w", att.Driver, err)
		}
		if err := d.Delete(att.URL); err != nil {
			return fmt.Errorf("删除文件失败: %s, 驱动: %s, err: %w", att.URL, att.Driver, err)
		}
	}

	// 删除数据库记录
	return s.IService.Delete(ctx, opts)
}
