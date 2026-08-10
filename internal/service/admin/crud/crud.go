package crud

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"
)

// TableInfo 数据表信息
type TableInfo struct {
	Table   string `json:"table"`
	Comment string `json:"comment"`
}

// FieldInfo 数据表字段信息
type FieldInfo struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	IsPK    bool   `json:"-"`
}

// ModelInfo 模型信息
type ModelInfo struct {
	Name    string `json:"name"`
	File    string `json:"file"`
	Package string `json:"package"`
	Comment string `json:"comment"`
}

// Service 可视化 CRUD 服务
type Service struct {
}

// NewService 创建可视化 CRUD 服务实例
func NewService() *Service {
	return &Service{}
}

// withPrefix 规范化表名: 已带前缀则先剥除再拼接，防止前缀重复
func withPrefix(name string) string {
	prefix := config.Get().Database.Prefix
	return prefix + strings.TrimPrefix(name, prefix)
}

// TableList 获取当前项目数据库中的数据表列表
func (s *Service) TableList(ctx context.Context, excludeTables []string) ([]TableInfo, error) {
	db := database.DB().WithContext(ctx)

	q := `
SELECT
    c.relname AS "table",
    COALESCE(obj_description(c.oid, 'pg_class'), '') AS "comment"
FROM pg_catalog.pg_class c
JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE c.relkind = 'r'
    AND n.nspname = current_schema()
`
	var args []any
	if len(excludeTables) > 0 {
		excluded := make([]string, 0, len(excludeTables))
		for _, name := range excludeTables {
			excluded = append(excluded, withPrefix(name))
		}
		q += "    AND c.relname NOT IN (?)\n"
		args = append(args, excluded)
	}
	q += "ORDER BY c.relname"

	var tables []TableInfo
	if err := db.Raw(q, args...).Scan(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

// TableExists 检查数据表是否存在
func (s *Service) TableExists(ctx context.Context, table string) (bool, error) {
	db := database.DB().WithContext(ctx)
	table = withPrefix(table)

	var count int64
	if err := db.Raw(`
SELECT COUNT(*) AS count
FROM pg_catalog.pg_class t
JOIN pg_catalog.pg_namespace n ON n.oid = t.relnamespace
WHERE n.nspname = current_schema()
    AND t.relkind = 'r'
    AND t.relname = ?`, table).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// TableFieldList 获取指定数据表的字段列表与主键
func (s *Service) TableFieldList(ctx context.Context, table string) (string, []FieldInfo, error) {
	db := database.DB().WithContext(ctx)
	table = withPrefix(table)

	var fields []FieldInfo
	fieldsSQL := `
SELECT
    a.attname AS "name",
    COALESCE(col_description(a.attrelid, a.attnum), '') AS "comment",
    EXISTS (
        SELECT 1
        FROM pg_catalog.pg_index i
        WHERE i.indrelid = a.attrelid
            AND i.indisprimary
            AND a.attnum = ANY(i.indkey)
    ) AS "is_pk"
FROM pg_catalog.pg_attribute a
JOIN pg_catalog.pg_class c ON c.oid = a.attrelid
JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = current_schema()
    AND c.relname = ?
    AND a.attnum > 0
    AND NOT a.attisdropped
ORDER BY a.attnum
`
	if err := db.Raw(fieldsSQL, table).Scan(&fields).Error; err != nil {
		return "", nil, err
	}

	// 遍历字段列表，找到主键字段
	pk := ""
	for _, f := range fields {
		if f.IsPK {
			pk = f.Name
			break
		}
	}
	return pk, fields, nil
}

// ModelList 获取模型列表
func (s *Service) ModelList() ([]ModelInfo, error) {
	files, err := filepath.Glob(filepath.Join(GenBaseDir["model"], "*.go"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)

	var models []ModelInfo
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}

		node, err := parser.ParseFile(token.NewFileSet(), file, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		for _, decl := range node.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				// 仅统计结构体类型
				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					continue
				}

				// 读取结构体上方注释，并去掉开头的结构体名称前缀
				comment := ""
				doc := typeSpec.Doc
				if doc == nil && len(genDecl.Specs) == 1 {
					doc = genDecl.Doc
				}
				if doc != nil {
					comment = strings.TrimSpace(doc.Text())
					comment = strings.TrimPrefix(comment, typeSpec.Name.Name+" ")
				}

				models = append(models, ModelInfo{
					Name:    typeSpec.Name.Name,
					Package: node.Name.Name,
					File:    filepath.ToSlash(file),
					Comment: comment,
				})
			}
		}
	}
	return models, nil
}
