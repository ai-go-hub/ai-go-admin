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
	"gorm.io/gorm/schema"
)

// Where 筛选条件
type Where struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

// WhereGroup 筛选条件分组
type WhereGroup struct {
	Wheres []Where `json:"wheres"` // 一组筛选条件
	Or     bool    `json:"or"`     // 组内条件是否使用 OR 连接，值为 false 则使用 AND 连接条件
}

// Options 通用服务操作选项
// Sort 以外的每个方法所需要的选项都可以在此找到，但并非每个方法都会使用全部的选项
type Options struct {
	OmitFields       []string             // 排除出入库字段，会传递给仓储层的 Omit 方法
	SelectFields     []string             // 选择出入库字段，会传递给仓储层的 Select 方法
	Wheres           []WhereGroup         // 查询条件组，用于构建 WhereScopes，然后传递给仓储层的 Scopes 方法
	SortField        string               // 排序字段，用于构建 OrderScope
	SortOrder        string               // 排序方式
	Page             int                  // 页码，用于构建 PaginateScope
	Limit            int                  // 每页条数
	Selector         bool                 // 是否为来自选择器的查询
	PrimaryKeyValue  string               // 主键值，目前可供 Get、Update 方法获取数据行
	PrimaryKeyValues []string             // 主键切片，目前可供 Delete 方法批量删除行
	Extension        any                  // 任意自定义扩展参数
	Preloads         []repository.Preload // 预加载关联
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
	BuildWhereScopes(wheres []WhereGroup) []func(*gorm.Statement)
	BuildWhereExpr(w Where) clause.Expression
	BuildRelationWhereExpr(sch *schema.Schema, segments []string, op string, value any) clause.Expression
	BuildLeafExpr(field, op string, value any) clause.Expression
	LookUpRelation(sch *schema.Schema, name string) *schema.Relationship
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
		Scopes: s.BuildWhereScopes(opts.Wheres),
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
	whereScopes := s.BuildWhereScopes(opts.Wheres)

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
		Preloads:         opts.Preloads,
		Scopes:           s.BuildScopes(opts),
		PrimaryKeyValue:  opts.PrimaryKeyValue,
		PrimaryKeyValues: opts.PrimaryKeyValues,
	}
}

// BuildScopes 统一组装过滤条件、排序、分页 Scopes
func (s *Service[T]) BuildScopes(opts Options) []func(*gorm.Statement) {
	// 构建 Where scopes
	scopes := s.BuildWhereScopes(opts.Wheres)

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
//
// 每个 WhereGroup 产出 1 个 scope
// 组内根据 Group.Or 字段决定 OR / AND 连接；组间恒定 AND
// 字段在模型中不存在的条件会被静默跳过
func (s *Service[T]) BuildWhereScopes(groups []WhereGroup) []func(*gorm.Statement) {
	if len(groups) == 0 {
		return nil
	}

	scopes := make([]func(*gorm.Statement), 0, len(groups))
	for _, g := range groups {
		if len(g.Wheres) == 0 {
			continue
		}

		exprs := make([]clause.Expression, 0, len(g.Wheres))
		for _, w := range g.Wheres {
			if expr := s.BuildWhereExpr(w); expr != nil {
				exprs = append(exprs, expr)
			}
		}
		if len(exprs) == 0 {
			continue
		}

		var wrapped clause.Expression
		if g.Or {
			wrapped = clause.OrConditions{Exprs: exprs}
		} else {
			wrapped = clause.AndConditions{Exprs: exprs}
		}
		scopes = append(scopes, func(stmt *gorm.Statement) {
			stmt.AddClause(clause.Where{Exprs: []clause.Expression{wrapped}})
		})
	}
	return scopes
}

// BuildWhereExpr 将单条 Where 转为 GORM clause.Expression
func (s *Service[T]) BuildWhereExpr(w Where) clause.Expression {
	field := strings.TrimSpace(w.Field)
	if field == "" {
		return nil
	}

	// 操作符号别名获取 + 白名单检查
	op := GetOperatorByAlias(w.Operator)
	if op == "" {
		return nil
	}

	// 关联字段名如: group.name，此处按 . 分段，从当前模型 Schema 出发递归下钻（支持嵌套关联）
	if strings.Contains(field, ".") {
		sch, err := s.repo.Schema()
		if err != nil {
			return nil
		}
		// 递归生成嵌套子查询
		return s.BuildRelationWhereExpr(sch, strings.Split(field, "."), op, w.Value)
	}

	// 本表字段
	if exists, _ := s.repo.FieldExists(field); !exists {
		return nil
	}
	return s.BuildLeafExpr(field, op, w.Value)
}

// BuildRelationWhereExpr 递归生成关联字段的嵌套子查询表达式
// segments 是按 . 切分后的路径段，如 ["group", "name"] 或 ["admin", "group", "name"]
func (s *Service[T]) BuildRelationWhereExpr(sch *schema.Schema, segments []string, op string, value any) clause.Expression {
	if len(segments) < 2 {
		return nil
	}

	// 使用双引号包裹标识符
	quote := func(parts ...string) string {
		if len(parts) == 0 {
			return ""
		}
		return `"` + strings.Join(parts, `"."`) + `"`
	}

	// 若 schema 存在名为 DeletedAt 的字段，返回追加到 WHERE 后的软删过滤 SQL 片段
	// 固定检查名为 DeletedAt 的字段，否则需要使用反射检查字段类型
	softDelete := func(s *schema.Schema) string {
		if s == nil {
			return ""
		}
		f := s.LookUpField("DeletedAt")
		if f == nil {
			return ""
		}
		return " AND " + quote(s.Table, f.DBName) + " IS NULL"
	}

	head, rest := segments[0], segments[1:]
	rel := s.LookUpRelation(sch, head)
	if rel == nil || len(rel.References) == 0 {
		return nil
	}

	relTable := rel.FieldSchema.Table

	// 生成关联表侧的 WHERE 子句: 可能是叶子字段，也可能是继续下钻
	var innerExpr clause.Expression
	if len(rest) == 1 {
		leafField := rel.FieldSchema.LookUpField(rest[0])
		if leafField == nil {
			return nil
		}
		innerExpr = s.BuildLeafExpr(quote(relTable, leafField.DBName), op, value)
	} else {
		innerExpr = s.BuildRelationWhereExpr(rel.FieldSchema, rest, op, value)
	}
	if innerExpr == nil {
		return nil
	}

	// 根据关联类型拼接 IN 子查询
	ref := rel.References[0]
	ownerTable := sch.Table
	relSoftDel := softDelete(rel.FieldSchema)

	switch rel.Type {

	// BelongsTo: 当前表.外键 IN (SELECT 关联表.主键 FROM 关联表 WHERE ...)
	// ForeignKey 在当前表，PrimaryKey 在关联表
	case schema.BelongsTo:
		return gorm.Expr(
			fmt.Sprintf("%s IN (SELECT %s FROM %s WHERE ?%s)",
				quote(ownerTable, ref.ForeignKey.DBName),
				quote(relTable, ref.PrimaryKey.DBName),
				quote(relTable),
				relSoftDel,
			), innerExpr)

	// HasOne/HasMany: 当前表.主键 IN (SELECT 关联表.外键 FROM 关联表 WHERE ...)
	// PrimaryKey 在当前表，ForeignKey 在关联表
	case schema.HasOne, schema.HasMany:
		return gorm.Expr(
			fmt.Sprintf("%s IN (SELECT %s FROM %s WHERE ?%s)",
				quote(ownerTable, ref.PrimaryKey.DBName),
				quote(relTable, ref.ForeignKey.DBName),
				quote(relTable),
				relSoftDel,
			), innerExpr)

	// Many2Many: 通过 JoinTable 中转两层 IN
	// References[0]: 当前表主键 <-> 中间表外键；References[1]: 关联表主键 <-> 中间表外键
	case schema.Many2Many:
		if rel.JoinTable == nil || len(rel.References) < 2 {
			return nil
		}
		joinTable := rel.JoinTable.Table
		var ownerRef, relRef *schema.Reference
		for _, r := range rel.References {
			if r.OwnPrimaryKey {
				ownerRef = r
			} else {
				relRef = r
			}
		}
		if ownerRef == nil || relRef == nil {
			return nil
		}
		joinSoftDel := softDelete(rel.JoinTable)
		return gorm.Expr(
			fmt.Sprintf("%s IN (SELECT %s FROM %s WHERE %s IN (SELECT %s FROM %s WHERE ?%s)%s)",
				quote(ownerTable, ownerRef.PrimaryKey.DBName),
				quote(joinTable, ownerRef.ForeignKey.DBName),
				quote(joinTable),
				quote(joinTable, relRef.ForeignKey.DBName),
				quote(relTable, relRef.PrimaryKey.DBName),
				quote(relTable),
				relSoftDel,
				joinSoftDel,
			), innerExpr)
	}

	return nil
}

// LookUpRelation 按名称查找关联，兼容 PascalCase 和 snake_case 两种命名风格
func (s *Service[T]) LookUpRelation(sch *schema.Schema, name string) *schema.Relationship {
	// 直接精确匹配
	if r, ok := sch.Relationships.Relations[name]; ok {
		return r
	}
	// 走 GORM 命名策略: 把请求名当作数据库列名，反向匹配 Go 字段名对应的列名
	namer := s.repo.DB().NamingStrategy
	target := namer.ColumnName(sch.Table, name)
	for k, r := range sch.Relationships.Relations {
		if namer.ColumnName(sch.Table, k) == target {
			return r
		}
	}
	return nil
}

// BuildLeafExpr 生成叶子（基础）的字段条件表达式
func (s *Service[T]) BuildLeafExpr(field, op string, value any) clause.Expression {
	switch op {
	case "IS NULL", "IS NOT NULL":
		return gorm.Expr(field + " " + op)
	case "BETWEEN", "NOT BETWEEN":
		str, ok := value.(string)
		if !ok {
			return nil
		}
		parts := strings.SplitN(str, ",", 2)
		if len(parts) != 2 {
			return nil
		}
		return gorm.Expr(field+" "+op+" ? AND ?", strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	case "LIKE", "NOT LIKE", "ILIKE", "NOT ILIKE":
		v := value
		if str, ok := v.(string); ok && !strings.Contains(str, "%") {
			v = "%" + str + "%"
		}
		return gorm.Expr(field+" "+op+" ?", v)
	case "IN", "NOT IN":
		return gorm.Expr(field+" "+op+" (?)", value)
	default:
		return gorm.Expr(field+" "+op+" ?", value)
	}
}

// GetOperatorByAlias 运算符别名 → SQL 运算符白名单
// 未在白名单内的返回空串，由调用方跳过该条件
func GetOperatorByAlias(op string) string {
	switch op {
	case "eq", "":
		return "="
	case "ne":
		return "!="
	case "gt":
		return ">"
	case "egt":
		return ">="
	case "lt":
		return "<"
	case "elt":
		return "<="
	case "LIKE", "NOT LIKE", "ILIKE", "NOT ILIKE":
		return op
	case "IN", "NOT IN":
		return op
	case "BETWEEN", "NOT BETWEEN":
		return op
	case "NULL":
		return "IS NULL"
	case "NOT NULL":
		return "IS NOT NULL"
	default:
		return ""
	}
}
