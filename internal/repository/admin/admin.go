package admin

import (
	"context"

	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"

	"gorm.io/gorm"
)

// AdminRepository 管理员仓储，嵌入通用仓储并扩展自定义方法
type AdminRepository struct {
	*repository.Repository[model.Admin]
}

// NewAdminRepository 创建管理员仓储实例
func NewAdminRepository() *AdminRepository {
	return &AdminRepository{
		Repository: repository.NewRepository[model.Admin](),
	}
}

// FindByUsername 根据用户名查询管理员
func (r *AdminRepository) FindByUsername(ctx context.Context, username string) (*model.Admin, error) {
	admin, err := gorm.G[model.Admin](r.DB()).Where("username = ?", username).First(ctx)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// UpdateLoginInfo 更新登录信息（登录时间、IP、重置失败次数）
func (r *AdminRepository) UpdateLoginInfo(ctx context.Context, id uint, loginIP string) error {
	// 注意：login_failure 置为 0 是零值，而基仓储的 Update 使用 struct 更新时，默认将跳过零值字段，
	// 所以需要单独使用 map[string]any 更新
	_, err := gorm.G[map[string]any](r.DB()).Table(model.Admin{}.TableName()).
		Where("id = ?", id).
		Updates(ctx, map[string]any{
			"last_login_ip": loginIP,
			"last_login_at": gorm.Expr("NOW()"),
			"login_failure": 0,
		})
	return err
}

// IncrementLoginFailure 增加登录失败次数
func (r *AdminRepository) IncrementLoginFailure(ctx context.Context, id uint) error {
	_, err := gorm.G[model.Admin](r.DB()).
		Where("id = ?", id).
		Update(ctx, "login_failure", gorm.Expr("login_failure + ?", 1))
	return err
}

// DeleteGroupAccesses 删除管理员的所有分组关联
func (r *AdminRepository) DeleteGroupAccesses(ctx context.Context, uid uint) error {
	_, err := gorm.G[model.AdminGroupAccess](r.DB()).Where("uid = ?", uid).Delete(ctx)
	return err
}

// ReplaceGroupAccesses 替换管理员的分组关联: 先删后插
func (r *AdminRepository) ReplaceGroupAccesses(ctx context.Context, uid uint, accesses []model.AdminGroupAccess) error {
	if err := r.DeleteGroupAccesses(ctx, uid); err != nil {
		return err
	}

	if len(accesses) == 0 {
		return nil
	}

	for i := range accesses {
		accesses[i].UID = uid
	}
	// 泛型 Create 只接受单条 *T，批量插入用传统 API
	return r.DB().WithContext(ctx).Create(&accesses).Error
}
