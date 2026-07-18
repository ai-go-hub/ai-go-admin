package repository

import (
	"errors"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Options 通用仓库操作选项（各选项可按需使用）
type Options struct {
	Omit       []string                // 排除出入库字段，将在 Create、Update、List 等方法中应用至 GORM 的 Select 方法
	Select     []string                // 选择出入库字段，其余同上
	Scopes     []func(*gorm.Statement) // 将在 List 和 Get 等方法中直接传递给 GORM 的 Scopes 方法，可以自定义 查询、排序、分页 等等
	PrimaryKey string                  // 主键值，提供则会在 Get 方法中作为 Where 的参数
}

// IRepository 通用仓库接口
type IRepository[T any] interface {
	Create(c *gin.Context, entity *T, opts Options) error
	Get(c *gin.Context, opts Options) (*T, error)
	List(c *gin.Context, opts Options) ([]T, error)
	Count(c *gin.Context, opts Options) (int64, error)
	Update(c *gin.Context, pk string, entity T, opts Options) error
	Delete(c *gin.Context, pk string) error
	PrimaryKeyField() (string, error)       // 获取主键字段
	FieldExists(field string) (bool, error) // 检查一个字段是否存在
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
func (r *Repository[T]) Create(c *gin.Context, entity *T, opts Options) error {
	var q gorm.CreateInterface[T] = gorm.G[T](r.DB())

	// 入库字段的选择与忽略
	if len(opts.Select) > 0 {
		q = q.Select(strings.Join(opts.Select, ","))
	}
	if len(opts.Omit) > 0 {
		q = q.Omit(opts.Omit...)
	}

	return q.Create(c.Request.Context(), entity)
}

// Get 查询单条记录
func (r *Repository[T]) Get(c *gin.Context, opts Options) (*T, error) {
	q := gorm.G[T](r.DB()).Scopes(opts.Scopes...)

	// 有提供主键
	if opts.PrimaryKey != "" {
		pk, err := r.PrimaryKeyField()
		if err != nil {
			return nil, err
		}
		q = q.Where(pk+" = ?", opts.PrimaryKey)
	}

	// 出库字段的选择与忽略
	if len(opts.Select) > 0 {
		q = q.Select(strings.Join(opts.Select, ","))
	}
	if len(opts.Omit) > 0 {
		q = q.Omit(opts.Omit...)
	}

	entity, err := q.First(c.Request.Context())
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// List 查询全部记录，排序、分页等可通过 opts.Scopes 传入
func (r *Repository[T]) List(c *gin.Context, opts Options) ([]T, error) {
	q := gorm.G[T](r.DB()).Scopes(opts.Scopes...)

	// 出库字段的选择与忽略
	if len(opts.Select) > 0 {
		q = q.Select(strings.Join(opts.Select, ","))
	}
	if len(opts.Omit) > 0 {
		q = q.Omit(opts.Omit...)
	}

	return q.Find(c.Request.Context())
}

// Count 统计记录总数
func (r *Repository[T]) Count(c *gin.Context, opts Options) (int64, error) {
	pkField, err := r.PrimaryKeyField()
	if err != nil {
		return 0, err
	}
	return gorm.G[T](r.DB()).Scopes(opts.Scopes...).Count(c.Request.Context(), pkField)
}

// Update 根据主键更新记录
func (r *Repository[T]) Update(c *gin.Context, pk string, entity T, opts Options) error {
	pkField, err := r.PrimaryKeyField()
	if err != nil {
		return err
	}

	q := gorm.G[T](r.DB()).Where(pkField+" = ?", pk)

	// 入库字段的选择与忽略
	if len(opts.Select) > 0 {
		q = q.Select(strings.Join(opts.Select, ","))
	}
	if len(opts.Omit) > 0 {
		q = q.Omit(opts.Omit...)
	}

	_, err = q.Updates(c.Request.Context(), entity)
	return err
}

// Delete 根据主键删除记录
func (r *Repository[T]) Delete(c *gin.Context, pk string) error {
	pkField, err := r.PrimaryKeyField()
	if err != nil {
		return err
	}
	_, err = gorm.G[T](r.DB()).Where(pkField+" = ?", pk).Delete(c.Request.Context())
	return err
}

// PrimaryKeyField 获取当前模型的主键数据库字段名
// 复合主键模型返回 GORM 识别的优先主键字段；Statement.Parse 会复用 GORM 的 Schema 缓存
func (r *Repository[T]) PrimaryKeyField() (string, error) {
	stmt := &gorm.Statement{DB: r.DB()}
	if err := stmt.Parse(new(T)); err != nil {
		return "", err
	}
	if stmt.Schema.PrioritizedPrimaryField == nil {
		return "", errors.New("模型未定义主键")
	}
	return stmt.Schema.PrioritizedPrimaryField.DBName, nil
}

// FieldExists 检查当前模型是否包含指定字段
// 同时支持 Go 字段名和数据库列名；Statement.Parse 会复用 GORM 的 Schema 缓存
func (r *Repository[T]) FieldExists(field string) (bool, error) {
	field = strings.TrimSpace(field)
	if field == "" {
		return false, nil
	}

	stmt := &gorm.Statement{DB: r.DB()}
	if err := stmt.Parse(new(T)); err != nil {
		return false, err
	}
	return stmt.Schema.LookUpField(field) != nil, nil
}
