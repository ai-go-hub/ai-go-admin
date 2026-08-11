package table

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"

	"gorm.io/gorm"
)

// TableExists 判断表是否存在
func TableExists(ctx context.Context, table string) bool {
	var exists bool
	if err := database.DB().WithContext(ctx).Raw("SELECT to_regclass(?) IS NOT NULL", table).Scan(&exists).Error; err != nil {
		return false
	}
	return exists
}

// DropTable 删除数据表
func DropTable(ctx context.Context, table string) error {
	return database.DB().WithContext(ctx).
		Exec("DROP TABLE IF EXISTS " + QuoteIdent(table)).Error
}

// ColumnExists 判断列是否存在
func ColumnExists(ctx context.Context, table, column string) bool {
	var count int64
	database.DB().WithContext(ctx).Raw(
		"SELECT COUNT(*) FROM information_schema.columns WHERE table_name = ? AND column_name = ?",
		table, column).Scan(&count)
	return count > 0
}

// CreateTable 按字段定义创建数据表
func CreateTable(ctx context.Context, tableName string, table dto.CRUDTable, fields []dto.CRUDFields, sqls *[]string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE TABLE IF NOT EXISTS %s (\n", QuoteIdent(tableName))
	cols := make([]string, 0, len(fields))
	for _, f := range fields {
		if f.Name == "" {
			continue
		}
		if f.PrimaryKey {
			cols = append(cols, "    "+QuoteIdent(f.Name)+" "+buildPKType(f)+" PRIMARY KEY")
		} else {
			cols = append(cols, "    "+QuoteIdent(f.Name)+" "+buildDDLColumn(f))
		}
	}
	b.WriteString(strings.Join(cols, ",\n"))
	b.WriteString("\n);")

	db := database.DB().WithContext(ctx)
	if err := execAndLog(db, sqls, b.String()); err != nil {
		return err
	}
	// 表注释
	if table.Comment != "" {
		commentSQL := "COMMENT ON TABLE " + QuoteIdent(tableName) + " IS '" + strings.ReplaceAll(table.Comment, "'", "''") + "'"
		if err := execAndLog(db, sqls, commentSQL); err != nil {
			return err
		}
	}
	// 列注释
	for _, f := range fields {
		if f.Comment != "" {
			if err := commentColumnDDL(db, tableName, f.Name, f.Comment, sqls); err != nil {
				return err
			}
		}
	}
	return nil
}

// SyncTableDesign 将 DesignChange 中的字段变更同步到已存在的数据表
func SyncTableDesign(ctx context.Context, tableName string, table dto.CRUDTable, fields []dto.CRUDFields, sqls *[]string) error {
	changes := table.DesignChange
	hasChange := false
	for _, c := range changes {
		if c.Sync {
			hasChange = true
			break
		}
	}
	if !hasChange {
		return nil
	}

	db := database.DB().WithContext(ctx)

	// 第一轮: 改名与删除（优先执行，避免先改名再改属性时找不到字段）
	for _, c := range changes {
		if !c.Sync {
			continue
		}
		switch c.Type {
		case "change-field-name":
			if !ColumnExists(ctx, tableName, c.Name) {
				return errors.New("重命名字段不存在: " + c.Name)
			}
			if err := execAndLog(db, sqls, fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", QuoteIdent(tableName), QuoteIdent(c.Name), QuoteIdent(c.NewName))); err != nil {
				return err
			}
		case "del-field":
			if f, ok := findFieldByName(fields, c.Name); ok && f.PrimaryKey {
				continue // 跳过主键字段删除
			}
			if !ColumnExists(ctx, tableName, c.Name) {
				return errors.New("删除字段不存在: " + c.Name)
			}
			if err := execAndLog(db, sqls, fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", QuoteIdent(tableName), QuoteIdent(c.Name))); err != nil {
				return err
			}
		}
	}

	// 第二轮: 修改属性与添加
	for _, c := range changes {
		if !c.Sync {
			continue
		}
		switch c.Type {
		case "change-field-attr":
			f, ok := findFieldByName(fields, c.Name)
			if !ok {
				return errors.New("修改属性字段不存在: " + c.Name)
			}
			if f.PrimaryKey {
				continue // 主键字段跳过
			}
			if !ColumnExists(ctx, tableName, c.Name) {
				return errors.New("修改属性字段不存在: " + c.Name)
			}
			if err := alterColumnDDL(ctx, db, tableName, f, sqls); err != nil {
				return err
			}
		case "add-field":
			f, ok := findFieldByName(fields, c.NewName)
			if !ok {
				return errors.New("添加字段不存在: " + c.NewName)
			}
			if f.PrimaryKey {
				continue // 主键字段跳过
			}
			if ColumnExists(ctx, tableName, c.NewName) {
				return errors.New("添加字段已存在: " + c.NewName)
			}
			if err := execAndLog(db, sqls, fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", QuoteIdent(tableName), QuoteIdent(f.Name), buildDDLColumn(f))); err != nil {
				return err
			}
			if err := commentColumnDDL(db, tableName, f.Name, f.Comment, sqls); err != nil {
				return err
			}
		}
	}
	return nil
}

// QuoteIdent 双引号包裹标识符
func QuoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

// findFieldByName 按字段名查找字段
func findFieldByName(fields []dto.CRUDFields, name string) (dto.CRUDFields, bool) {
	for _, f := range fields {
		if f.Name == name {
			return f, true
		}
	}
	return dto.CRUDFields{}, false
}

// buildDDLType 生成列类型
func buildDDLType(f dto.CRUDFields) string {
	switch f.Type {
	case "varchar", "char":
		if f.Length > 0 {
			return fmt.Sprintf("%s(%d)", f.Type, f.Length)
		}
	case "numeric":
		if f.Length > 0 && f.Precision > 0 {
			return fmt.Sprintf("numeric(%d,%d)", f.Length, f.Precision)
		}
	}
	return f.Type
}

// buildDDLDefault 生成默认值表达式；无默认值时返回空串
func buildDDLDefault(f dto.CRUDFields) string {
	switch f.DefaultType {
	case "EMPTY STRING":
		return "''"
	case "INPUT":
		if f.Default == "" {
			return ""
		}
		switch f.Type {
		case "jsonb":
			return "'" + jsonbDefault(f.Default) + "'::jsonb"
		case "varchar", "char", "text", "uuid":
			return "'" + strings.ReplaceAll(f.Default, "'", "''") + "'"
		default:
			return f.Default
		}
	}
	return ""
}

// jsonbDefault 将 jsonb 默认值规范为合法 JSON 数组文本
func jsonbDefault(v string) string {
	v = strings.TrimSpace(v)
	if strings.HasPrefix(v, "[") {
		// 已组装为数组: 单引号数组规范为双引号，合法 JSON 原样保留
		if !json.Valid([]byte(v)) {
			return strings.ReplaceAll(v, "'", `"`)
		}
		return v
	}
	// 逗号分隔形式转换 opt0,opt1 -> ["opt0","opt1"]
	parts := strings.Split(v, ",")
	for i := range parts {
		parts[i] = `"` + strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(parts[i]), `\`, `\\`), `"`, `\"`) + `"`
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// buildDDLColumn 生成完整列定义（类型 + NOT NULL + DEFAULT）
func buildDDLColumn(f dto.CRUDFields) string {
	var b strings.Builder
	b.WriteString(buildDDLType(f))
	if !f.Null && !f.PrimaryKey {
		b.WriteString(" NOT NULL")
	}
	if f.Unique && !f.PrimaryKey {
		b.WriteString(" UNIQUE")
	}
	if def := buildDDLDefault(f); def != "" {
		fmt.Fprintf(&b, " DEFAULT %s", def)
	}
	return b.String()
}

// buildPKType 生成主键列类型
func buildPKType(f dto.CRUDFields) string {
	switch f.Generated {
	case "GENERATED ALWAYS":
		return buildDDLType(f) + " GENERATED ALWAYS AS IDENTITY"
	case "GENERATED BY DEFAULT":
		return buildDDLType(f) + " GENERATED BY DEFAULT AS IDENTITY"
	}
	return buildDDLType(f)
}

// alterColumnDDL 修改列属性（类型/可空/默认值/注释/唯一约束）
func alterColumnDDL(ctx context.Context, db *gorm.DB, table string, f dto.CRUDFields, sqls *[]string) error {
	col := QuoteIdent(f.Name)
	// 类型
	if err := execAndLog(db, sqls, fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", QuoteIdent(table), col, buildDDLType(f))); err != nil {
		return err
	}
	// 可空
	nullOp := "DROP NOT NULL"
	if !f.Null && !f.PrimaryKey {
		nullOp = "SET NOT NULL"
	}
	if err := execAndLog(db, sqls, fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s %s", QuoteIdent(table), col, nullOp)); err != nil {
		return err
	}
	// 默认值
	if def := buildDDLDefault(f); def != "" {
		if err := execAndLog(db, sqls, fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET DEFAULT %s", QuoteIdent(table), col, def)); err != nil {
			return err
		}
	} else if err := execAndLog(db, sqls, fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s DROP DEFAULT", QuoteIdent(table), col)); err != nil {
		return err
	}
	// 唯一约束
	if err := alterUniqueConstraint(ctx, db, table, f, sqls); err != nil {
		return err
	}
	// 注释
	return commentColumnDDL(db, table, f.Name, f.Comment, sqls)
}

// alterUniqueConstraint 按字段 Unique 属性添加/删除唯一约束（主键隐含唯一，跳过）
func alterUniqueConstraint(ctx context.Context, db *gorm.DB, tableName string, f dto.CRUDFields, sqls *[]string) error {
	if f.PrimaryKey {
		return nil
	}
	existing := columnUniqueConstraints(ctx, tableName, f.Name)
	if f.Unique {
		if len(existing) == 0 {
			constraint := fmt.Sprintf("%s_%s_key", tableName, f.Name)
			return execAndLog(db, sqls, fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s)", QuoteIdent(tableName), QuoteIdent(constraint), QuoteIdent(f.Name)))
		}
		return nil
	}
	for _, name := range existing {
		if err := execAndLog(db, sqls, fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s", QuoteIdent(tableName), QuoteIdent(name))); err != nil {
			return err
		}
	}
	return nil
}

// commentColumnDDL 设置列注释（COMMENT 语句不接受参数占位符，需内联并转义单引号）
func commentColumnDDL(db *gorm.DB, table, column, comment string, sqls *[]string) error {
	if comment == "" {
		return nil
	}
	sql := "COMMENT ON COLUMN " + QuoteIdent(table) + "." + QuoteIdent(column) + " IS '" +
		strings.ReplaceAll(comment, "'", "''") + "'"
	return execAndLog(db, sqls, sql)
}

// columnUniqueConstraints 查询列上的唯一约束名列表（排除主键）
func columnUniqueConstraints(ctx context.Context, tableName, column string) []string {
	var names []string
	database.DB().WithContext(ctx).Raw(`
SELECT tc.constraint_name
FROM information_schema.table_constraints tc
JOIN information_schema.key_column_usage kcu
  ON kcu.constraint_name = tc.constraint_name AND kcu.table_schema = tc.table_schema
WHERE tc.table_name = ? AND tc.constraint_type = 'UNIQUE' AND kcu.column_name = ?`,
		tableName, column).Scan(&names)
	return names
}

// execAndLog 执行 SQL 并追加到日志列表
func execAndLog(db *gorm.DB, sqls *[]string, sql string) error {
	if sqls != nil {
		*sqls = append(*sqls, sql)
	}
	return db.Exec(sql).Error
}
