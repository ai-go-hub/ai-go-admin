package admin

import (
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"

	"github.com/gin-gonic/gin"
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
func (r *AdminRepository) FindByUsername(c *gin.Context, username string) (*model.Admin, error) {
	admin, err := gorm.G[model.Admin](r.DB()).Where("username = ?", username).First(c.Request.Context())
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// UpdateLoginInfo 更新登录信息（登录时间、IP、重置失败次数）
func (r *AdminRepository) UpdateLoginInfo(c *gin.Context, id uint, loginIP string) error {
	// 注意：login_failure 置为 0 是零值，而基仓储的 Update 使用 struct 更新时，默认将跳过零值字段，
	// 所以需要单独使用 map[string]any 更新
	_, err := gorm.G[map[string]any](r.DB()).Table(model.Admin{}.TableName()).
		Where("id = ?", id).
		Updates(c.Request.Context(), map[string]any{
			"last_login_ip": loginIP,
			"last_login_at": gorm.Expr("NOW()"),
			"login_failure": 0,
		})
	return err
}

// IncrementLoginFailure 增加登录失败次数
func (r *AdminRepository) IncrementLoginFailure(c *gin.Context, id uint) error {
	_, err := gorm.G[model.Admin](r.DB()).
		Where("id = ?", id).
		Update(c.Request.Context(), "login_failure", gorm.Expr("login_failure + ?", 1))
	return err
}
