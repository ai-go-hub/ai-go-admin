package common

import (
	"context"
	"errors"
	"time"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"

	"gorm.io/gorm"
)

// AttachmentRepository 附件仓储，嵌入通用仓储并扩展附件专属查询
type AttachmentRepository struct {
	*repository.Repository[model.Attachment]
}

// NewAttachmentRepository 创建附件仓储实例
func NewAttachmentRepository() *AttachmentRepository {
	return &AttachmentRepository{
		Repository: repository.NewRepository[model.Attachment](),
	}
}

// FindBySha1 根据 SHA1 查询附件，未找到返回 nil
func (r *AttachmentRepository) FindBySha1(ctx context.Context, sha1 string) (*model.Attachment, error) {
	att, err := gorm.G[model.Attachment](r.DB()).Where("sha1 = ?", sha1).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &att, nil
}

// IncrementQuote 附件引用计数 +1，并更新最后上传时间与上传者
func (r *AttachmentRepository) IncrementQuote(ctx context.Context, id uint, now time.Time, userID uint, userType string) error {
	updates := map[string]any{
		"quote":          gorm.Expr("quote + ?", 1),
		"last_upload_at": now,
		"user_id":        userID,
		"user_type":      userType,
	}
	_, err := gorm.G[map[string]any](r.DB()).
		Table(model.Attachment{}.TableName()).
		Where("id = ?", id).
		Updates(ctx, updates)
	return err
}

// Save 入库附件记录
func (r *AttachmentRepository) Save(ctx context.Context, att *model.Attachment) error {
	return gorm.G[model.Attachment](r.DB()).Create(ctx, att)
}
