package service

import (
	"errors"
	"fmt"
	"slices"
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
	Sort(c *gin.Context, opts Options, move, target, direction string, weigh int64) error
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

// Sort 排序，
// 使用 `增量重排法` 或叫 `区间位移法`（而不是交换法）
func (s *Service[T]) Sort(c *gin.Context, opts Options, move, target, direction string, weigh int64) error {
	pkField, err := s.repo.PrimaryKeyField()
	if err != nil {
		return err
	}

	weighField := strings.TrimSpace(opts.SortField)
	if weighField == "" {
		return errors.New("缺少排序字段")
	}
	if exists, _ := s.repo.FieldExists(weighField); !exists {
		return errors.New("排序字段错误")
	}

	move = strings.TrimSpace(move)
	target = strings.TrimSpace(target)
	if move == "" || target == "" || move == target {
		return errors.New("移动行和目标行必填且不能相同")
	}

	ctx := c.Request.Context()
	db := s.repo.DB().WithContext(ctx)

	// 波及行权重变化方向
	// 用户排序方式 asc （权重值小的在前）: 向上拖 → 波及行权重 +1；向下拖 → 波及行权重 -1
	// 用户排序方式 desc（权重值大的在前）: 向上拖 → 波及行权重 -1；向下拖 → 波及行权重 +1
	sortOrderDesc := strings.EqualFold(strings.TrimSpace(opts.SortOrder), "desc")

	updateOp := "-"
	if (sortOrderDesc && direction == "down") || (!sortOrderDesc && direction == "up") {
		updateOp = "+"
	}

	// 复用 BuildOrderScope，保证与列表排序规则一致
	orderScope := s.BuildOrderScope(opts.SortField, opts.SortOrder)

	// 构建过滤条件 scopes（与列表查询共享同一套筛选逻辑）
	whereScopes := BuildWhereScopes(opts.Wheres)

	// 将 func(*gorm.Statement) scopes 转为传统 API 的 func(*gorm.DB) 格式
	// 后续需要使用 GORM 传统 API 的 Pluck 方法；不使用它则需要使用反射获取权重字段值等，更为复杂
	stmtScopes := make([]func(*gorm.DB) *gorm.DB, 0, len(whereScopes)+1)
	for _, scope := range whereScopes {
		s := scope
		stmtScopes = append(stmtScopes, func(tx *gorm.DB) *gorm.DB { s(tx.Statement); return tx })
	}
	if orderScope != nil {
		stmtScopes = append(stmtScopes, func(tx *gorm.DB) *gorm.DB { orderScope(tx.Statement); return tx })
	}

	// 查询与目标行同权重的所有行主键
	var weighIDsAny []any
	if err := db.Model(new(T)).Scopes(stmtScopes...).Where(weighField+" = ?", weigh).Pluck(pkField, &weighIDsAny).Error; err != nil {
		return err
	}

	// weighIDs 统一转字符串以便后续比较（避免 uint/int64/string 混杂时的类型不匹配）
	weighIDs := make([]string, len(weighIDsAny))
	for i, v := range weighIDsAny {
		weighIDs[i] = fmt.Sprint(v)
	}
	weighRowsCount := len(weighIDs)

	// 事务: 批量位移 + 单行增量重排
	return db.Transaction(func(tx *gorm.DB) error {
		gtx := gorm.G[T](tx).Scopes(whereScopes...)

		// 一次 SQL 完成波及行的整体挪位
		// dec（权重 −）: 波及 weigh < target 的行 → 再减，向下移位
		// inc（权重 +）: 波及 weigh > target 的行 → 再加，向上移位
		bulkOp := "<"
		if updateOp == "+" {
			bulkOp = ">"
		}

		if _, err := gtx.Where(weighField+" "+bulkOp+" ?", weigh).Where(pkField+" <> ?", move).
			Update(ctx, weighField, gorm.Expr(weighField+" "+updateOp+" ?", weighRowsCount)); err != nil {
			return err
		}

		// 向下拖动时反转，保证等权重区间内相对顺序不变
		if direction == "down" {
			slices.Reverse(weighIDs)
		}

		exprOffset := func(offset int) int64 {
			if updateOp == "+" {
				return weigh + int64(offset)
			}
			return weigh - int64(offset)
		}

		// 遍历等权重行，每匹配到一行时，权重再额外挪 1 位
		moveComplete := 0
		for i, rowID := range weighIDs {
			// 跳过被拖动行本身（等权重区间内互拖时会出现）
			if rowID == move {
				continue
			}

			// 当前行相对目标权重的偏移
			rowWeigh := exprOffset(i)

			// 命中目标行: 将被拖动行放到此处
			if rowID == target {
				if _, err := gtx.Where(pkField+" = ?", move).Update(ctx, weighField, rowWeigh); err != nil {
					return err
				}
				moveComplete = 1
			}

			// 目标行命中后，剩余等权重行的偏移额外 +1（腾出被拖动行的位置）
			if moveComplete == 1 {
				rowWeigh = exprOffset(i + moveComplete)
			}

			if _, err := gtx.Where(pkField+" = ?", rowID).Update(ctx, weighField, rowWeigh); err != nil {
				return err
			}
		}
		return nil
	})
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
	if s := s.BuildOrderScope(opts.SortField, opts.SortOrder); s != nil {
		scopes = append(scopes, s)
	}

	// 构建 Paginate scope
	if s := BuildPaginateScope(opts.Page, opts.Limit); s != nil {
		scopes = append(scopes, s)
	}
	return scopes
}

// BuildOrderScope 构建排序 Scope
// 用户排序优先，主键 DESC 托底作为有序保证，确保分页结果稳定
func (s *Service[T]) BuildOrderScope(sortField, sortOrder string) func(*gorm.Statement) {
	var clauses []string

	// 用户排序优先（校验字段是否存在，不存在的字段静默跳过）
	if field := strings.TrimSpace(sortField); field != "" {
		if exists, _ := s.repo.FieldExists(field); exists {
			switch strings.ToLower(strings.TrimSpace(sortOrder)) {
			case "desc":
				clauses = append(clauses, field+" DESC")
			default:
				clauses = append(clauses, field+" ASC")
			}
		}
	}

	// 主键托底排序，避免分页结果因排序不稳定而重复或遗漏
	// 数据库一般不保证排序，即先查到哪行就输出哪行且不保证多次相同查询中的输出顺序
	pkField, _ := s.repo.PrimaryKeyField()
	if pkField != "" {
		clauses = append(clauses, pkField+" DESC")
	}

	if len(clauses) == 0 {
		return nil
	}

	return func(stmt *gorm.Statement) {
		for _, c := range clauses {
			stmt.Order(c)
		}
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
