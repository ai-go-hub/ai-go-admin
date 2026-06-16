package repository

import (
	"ai-go-mall/internal/infra/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IRepository 通用仓库接口
type IRepository[T any] interface {
	Create(c *gin.Context, entity *T) error
	Get(c *gin.Context, id uint) (*T, error)
	List(c *gin.Context, scopes ...func(*gorm.Statement)) ([]T, error)
	Update(c *gin.Context, id uint, entity T) error
	Delete(c *gin.Context, id uint) error
}

// Repository 通用仓库实现
type Repository[T any] struct {
	db *gorm.DB
}

// RepositoryConfig 仓库配置选项
type RepositoryConfig struct {
	db *gorm.DB
}

// Option 选项函数
type Option func(*RepositoryConfig)

// WithDB 自定义 DB 实例
func WithDB(db *gorm.DB) Option {
	return func(c *RepositoryConfig) { c.db = db }
}

// NewRepository 创建通用仓库实例，支持函数式选项
func NewRepository[T any](opts ...Option) *Repository[T] {
	cfg := &RepositoryConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return &Repository[T]{db: cfg.db}
}

// WithDB 以链式调用的方式自定义 DB 实例
func (r *Repository[T]) WithDB(db *gorm.DB) *Repository[T] {
	r.db = db
	return r
}

// DB 获取当前使用的 DB 实例，优先返回自定义实例，否则返回全局实例
func (r *Repository[T]) DB() *gorm.DB {
	if r.db != nil {
		return r.db
	}
	return database.DB()
}

// Create 创建记录
func (r *Repository[T]) Create(c *gin.Context, entity *T) error {
	return gorm.G[T](r.DB()).Create(c.Request.Context(), entity)
}

// Get 根据主键 ID 查询单条记录
func (r *Repository[T]) Get(c *gin.Context, id uint) (*T, error) {
	entity, err := gorm.G[T](r.DB()).Where("id = ?", id).First(c.Request.Context())
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// List 查询全部记录，支持 GORM Scopes 动态条件
func (r *Repository[T]) List(c *gin.Context, scopes ...func(*gorm.Statement)) ([]T, error) {
	return gorm.G[T](r.DB()).Scopes(scopes...).Find(c.Request.Context())
}

// Update 根据主键 ID 更新记录
func (r *Repository[T]) Update(c *gin.Context, id uint, entity T) error {
	_, err := gorm.G[T](r.DB()).Where("id = ?", id).Updates(c.Request.Context(), entity)
	return err
}

// Delete 根据主键 ID 删除记录
func (r *Repository[T]) Delete(c *gin.Context, id uint) error {
	_, err := gorm.G[T](r.DB()).Where("id = ?", id).Delete(c.Request.Context())
	return err
}
