package service

import (
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Where 筛选条件
type Where struct {
	Field    string
	Operator string
	Value    any
}

// Options 通用服务操作选项
// 每个方法所需要的选项都可以在此找到，但并非每个方法都会使用全部的选项
type Options struct {
	OmitFields       []string // 排除出入库字段，会传递给仓储层的 Omit 方法
	SelectFields     []string // 选择出入库字段，会传递给仓储层的 Select 方法
	Wheres           []Where  // 查询条件，用于构建 WhereScopes，然后传递给仓储层的 Scopes 方法
	SortField        string   // 排序字段，用于构建 OrderScope
	SortOrder        string   // 排序方式
	Page             int      // 页码，用于构建 PaginateScope
	Limit            int      // 每页条数
	PrimaryKeyValue  string   // 主键值，目前可供 Get、Update 方法获取数据行
	PrimaryKeyValues []string // 主键切片，目前可供 Delete 方法批量删除行
	Extension        any      // 任意自定义扩展参数
}

// IService 通用服务接口
type IService[T any] interface {
	Create(c *gin.Context, entity *T, opts Options) error
	Get(c *gin.Context, opts Options) (*T, error)
	List(c *gin.Context, opts Options) ([]T, error)
	Count(c *gin.Context, opts Options) (int64, error)
	Update(c *gin.Context, entity *T, opts Options) error
	Delete(c *gin.Context, opts Options) error
	BuildRepoOpts(opts Options) repository.Options
	BuildScopes(opts Options) []func(*gorm.Statement)
}

// Service 通用服务实现
type Service[T any] struct {
	repo repository.IRepository[T]
}

// NewService 创建通用服务实例
func NewService[T any](repo repository.IRepository[T]) IService[T] {
	return &Service[T]{repo: repo}
}

// Create 创建记录
func (s *Service[T]) Create(c *gin.Context, entity *T, opts Options) error {
	return s.repo.Create(c, entity, s.BuildRepoOpts(opts))
}

// Get 查询单条记录
func (s *Service[T]) Get(c *gin.Context, opts Options) (*T, error) {
	return s.repo.Get(c, s.BuildRepoOpts(opts))
}

// List 查询全部记录
func (s *Service[T]) List(c *gin.Context, opts Options) ([]T, error) {
	return s.repo.List(c, s.BuildRepoOpts(opts))
}

// Count 统计满足过滤条件的记录总数
func (s *Service[T]) Count(c *gin.Context, opts Options) (int64, error) {
	// 只传递 Where scopes，忽略排序、分页等
	return s.repo.Count(c, repository.Options{
		Scopes: BuildWhereScopes(opts.Wheres),
	})
}

// Update 根据主键更新记录
func (s *Service[T]) Update(c *gin.Context, entity *T, opts Options) error {
	return s.repo.Update(c, entity, s.BuildRepoOpts(opts))
}

// Delete 根据主键批量删除记录
func (s *Service[T]) Delete(c *gin.Context, opts Options) error {
	return s.repo.Delete(c, s.BuildRepoOpts(opts))
}

// BuildRepoOpts 构建仓储层选项数据
func (s *Service[T]) BuildRepoOpts(opts Options) repository.Options {
	return repository.Options{
		OmitFields:       opts.OmitFields,
		SelectFields:     opts.SelectFields,
		Scopes:           s.BuildScopes(opts),
		PrimaryKeyValue:  opts.PrimaryKeyValue,
		PrimaryKeyValues: opts.PrimaryKeyValues,
	}
}

// BuildScopes 统一组装过滤条件、排序、分页 Scopes
func (s *Service[T]) BuildScopes(opts Options) []func(*gorm.Statement) {
	// 构建 Where scopes
	scopes := BuildWhereScopes(opts.Wheres)

	// 构建 Order scope
	// 排序字段使用 FieldExists 校验，不存在的字段静默不排序
	if opts.SortField != "" {
		if exists, _ := s.repo.FieldExists(opts.SortField); exists {
			if s := BuildOrderScope(opts.SortField, opts.SortOrder); s != nil {
				scopes = append(scopes, s)
			}
		}
	}

	// 构建 Paginate scope
	if s := BuildPaginateScope(opts.Page, opts.Limit); s != nil {
		scopes = append(scopes, s)
	}
	return scopes
}

// BuildOrderScope 构建排序 Scope
func BuildOrderScope(sortField, sortOrder string) func(*gorm.Statement) {
	field := strings.TrimSpace(sortField)
	if field == "" {
		return nil
	}
	switch strings.ToLower(strings.TrimSpace(sortOrder)) {
	case "desc":
		field += " DESC"
	default:
		field += " ASC"
	}
	return func(stmt *gorm.Statement) {
		stmt.Order(field)
	}
}

// BuildPaginateScope 构建分页 Scope
// limit > 200 截为 200
func BuildPaginateScope(page, limit int) func(*gorm.Statement) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return func(stmt *gorm.Statement) {
		stmt.AddClause(clause.Limit{Limit: &limit, Offset: offset})
	}
}

// BuildWhereScopes 构建查询条件 scopes
func BuildWhereScopes(wheres []Where) []func(*gorm.Statement) {
	if len(wheres) == 0 {
		return nil
	}
	scopes := make([]func(*gorm.Statement), 0, len(wheres))
	for _, w := range wheres {
		op := GetOperatorByAlias(w.Operator)
		switch op {
		case "IS NULL", "IS NOT NULL":
			scopes = append(scopes, func(stmt *gorm.Statement) {
				stmt.Where(w.Field + " " + op)
			})
		case "BETWEEN":
			scopes = append(scopes, func(stmt *gorm.Statement) {
				stmt.Where(w.Field+" BETWEEN ? AND ?", w.Value)
			})
		case "LIKE", "NOT LIKE", "ILIKE", "NOT ILIKE":
			v := w.Value
			if s, ok := v.(string); ok && !strings.Contains(s, "%") {
				v = "%" + s + "%"
			}
			scopes = append(scopes, func(stmt *gorm.Statement) {
				stmt.Where(w.Field+" "+w.Operator+" ?", v)
			})
		default:
			scopes = append(scopes, func(stmt *gorm.Statement) {
				stmt.Where(w.Field+" "+op+" ?", w.Value)
			})
		}
	}
	return scopes
}

// GetOperatorByAlias 符号类运算符别名 → SQL 运算符
func GetOperatorByAlias(op string) string {
	switch op {
	case "eq", "":
		return "="
	case "ne":
		return "!="
	case "gt":
		return ">"
	case "gte":
		return ">="
	case "lt":
		return "<"
	case "lte":
		return "<="
	default:
		return op // LIKE、IN、NOT IN 等单词运算符直接透传
	}
}
