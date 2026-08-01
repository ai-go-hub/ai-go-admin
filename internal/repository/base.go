package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Preload 预加载关联配置
type Preload struct {
	Association string                          // 关联名称
	Query       func(gorm.PreloadBuilder) error // 可选的自定义查询条件，为 nil 时仅按名称预加载
}

// Options 通用仓库操作选项
// 每个方法所需要的选项都可以在此找到，但并非每个方法都会使用全部的选项
type Options struct {
	Scopes           []func(*gorm.Statement) // 将在 Get 和 List 等方法中直接传递给 GORM 的 Scopes 方法，可以自定义 查询、排序、分页 等等
	OmitFields       []string                // 排除出入库字段，将在 Create、Update、List 等方法中应用至 GORM 的 Select 方法
	SelectFields     []string                // 选择出入库字段，其余同上
	PrimaryKeyValue  string                  // 主键值，目前可辅助 Get、Update 方法获取数据行
	PrimaryKeyValues []string                // 主键切片，目前可供 Delete 方法批量删除行
	Preloads         []Preload               // 预加载关联，将在 Get、List 中传递给 GORM 的 Preload 方法
}

// IRepository 通用仓库接口
type IRepository[T any] interface {
	DB() *gorm.DB
	Create(ctx context.Context, entity *T, opts Options) error
	Get(ctx context.Context, opts Options) (*T, error)
	List(ctx context.Context, opts Options) ([]T, error)
	Count(ctx context.Context, opts Options) (int64, error)
	Update(ctx context.Context, entity *T, opts Options) error
	Delete(ctx context.Context, opts Options) error
	PrimaryKeyField() (string, error)       // 获取主键字段
	FieldExists(field string) (bool, error) // 检查一个字段是否存在
	Schema() (*schema.Schema, error)        // 获取当前模型的 GORM Schema
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
func (r *Repository[T]) Create(ctx context.Context, entity *T, opts Options) error {
	var q gorm.CreateInterface[T] = gorm.G[T](r.DB())

	// 入库字段的选择
	if len(opts.SelectFields) > 0 {
		q = q.Select(strings.Join(opts.SelectFields, ","))
	}

	// 入库字段的忽略，未设置则忽略主键字段
	omitFields := opts.OmitFields
	if len(omitFields) == 0 {
		pk, _ := r.PrimaryKeyField()
		if pk != "" {
			omitFields = append(omitFields, pk)
		}
	}
	if len(omitFields) > 0 {
		q = q.Omit(omitFields...)
	}

	return q.Create(ctx, entity)
}

// Get 查询单条记录
func (r *Repository[T]) Get(ctx context.Context, opts Options) (*T, error) {
	q := gorm.G[T](r.DB()).Scopes(opts.Scopes...)

	// 预加载关联
	for _, p := range opts.Preloads {
		q = q.Preload(p.Association, p.Query)
	}

	// 有提供主键
	if opts.PrimaryKeyValue != "" {
		pk, err := r.PrimaryKeyField()
		if err != nil {
			return nil, err
		}
		q = q.Where(pk+" = ?", opts.PrimaryKeyValue)
	}

	// 出库字段的选择与忽略
	if len(opts.SelectFields) > 0 {
		q = q.Select(strings.Join(opts.SelectFields, ","))
	}
	if len(opts.OmitFields) > 0 {
		q = q.Omit(opts.OmitFields...)
	}

	entity, err := q.First(ctx)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// List 查询全部记录，排序、分页等可通过 opts.Scopes 传入
func (r *Repository[T]) List(ctx context.Context, opts Options) ([]T, error) {
	q := gorm.G[T](r.DB()).Scopes(opts.Scopes...)

	// 预加载关联
	for _, p := range opts.Preloads {
		q = q.Preload(p.Association, p.Query)
	}

	// 出库字段的选择与忽略
	if len(opts.SelectFields) > 0 {
		q = q.Select(strings.Join(opts.SelectFields, ","))
	}
	if len(opts.OmitFields) > 0 {
		q = q.Omit(opts.OmitFields...)
	}

	return q.Find(ctx)
}

// Count 统计记录总数
func (r *Repository[T]) Count(ctx context.Context, opts Options) (int64, error) {
	pkField, err := r.PrimaryKeyField()
	if err != nil {
		return 0, err
	}
	return gorm.G[T](r.DB()).Scopes(opts.Scopes...).Count(ctx, pkField)
}

// Update 根据主键更新记录
func (r *Repository[T]) Update(ctx context.Context, entity *T, opts Options) error {
	pkField, err := r.PrimaryKeyField()
	if err != nil {
		return err
	}

	q := gorm.G[T](r.DB()).Where(pkField+" = ?", opts.PrimaryKeyValue)

	// 入库字段的选择与忽略
	if len(opts.SelectFields) > 0 {
		q = q.Select(strings.Join(opts.SelectFields, ","))
	}
	if len(opts.OmitFields) > 0 {
		q = q.Omit(opts.OmitFields...)
	}

	_, err = q.Updates(ctx, *entity)
	return err
}

// Delete 根据主键批量删除记录
func (r *Repository[T]) Delete(ctx context.Context, opts Options) error {
	pkField, err := r.PrimaryKeyField()
	if err != nil {
		return err
	}
	_, err = gorm.G[T](r.DB()).Where(pkField+" IN ?", opts.PrimaryKeyValues).Delete(ctx)
	return err
}

// PrimaryKeyField 获取当前模型的主键数据库字段名
// 复合主键模型返回 GORM 识别的优先主键字段
func (r *Repository[T]) PrimaryKeyField() (string, error) {
	sch, err := r.Schema()
	if err != nil {
		return "", err
	}
	if sch.PrioritizedPrimaryField == nil {
		return "", errors.New("模型未定义主键")
	}
	return sch.PrioritizedPrimaryField.DBName, nil
}

// FieldExists 检查当前模型是否包含指定字段
// 同时支持 Go 字段名和数据库列名
func (r *Repository[T]) FieldExists(field string) (bool, error) {
	field = strings.TrimSpace(field)
	if field == "" {
		return false, nil
	}
	sch, err := r.Schema()
	if err != nil {
		return false, err
	}
	return sch.LookUpField(field) != nil, nil
}

// Schema 获取当前模型解析后的 GORM Schema
// Statement.Parse 会复用 GORM 的 Schema 缓存
func (r *Repository[T]) Schema() (*schema.Schema, error) {
	stmt := &gorm.Statement{DB: r.DB()}
	if err := stmt.Parse(new(T)); err != nil {
		return nil, err
	}
	return stmt.Schema, nil
}
