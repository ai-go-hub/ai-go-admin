package crud

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	repoAuth "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/auth"
	repoCommon "github.com/ai-go-hub/ai-go-admin/internal/repository/common"
	tbl "github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud/table"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"
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
	configRepo *repoCommon.ConfigRepository
}

// NewService 创建可视化 CRUD 服务实例
func NewService(configRepo *repoCommon.ConfigRepository) *Service {
	return &Service{
		configRepo: configRepo,
	}
}

// TableList 获取当前项目数据库中的数据表列表
func (s *Service) TableList(ctx context.Context, exclusions []string) ([]TableInfo, error) {
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
	if len(exclusions) > 0 {
		excluded := make([]string, 0, len(exclusions))
		for _, name := range exclusions {
			excluded = append(excluded, WithPrefix(name))
		}
		q += "    AND c.relname NOT IN (?)\n"
		args = append(args, excluded)
	}
	q += "ORDER BY c.relname"

	var tables []TableInfo
	if err := db.Raw(q, args...).Scan(&tables).Error; err != nil {
		return nil, err
	}
	sort.Slice(tables, func(i, j int) bool { return tables[i].Table < tables[j].Table })
	return tables, nil
}

// TableFieldList 获取指定数据表的字段列表与主键
func (s *Service) TableFieldList(ctx context.Context, table string) (string, []FieldInfo, error) {
	db := database.DB().WithContext(ctx)
	table = WithPrefix(table)

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
func (s *Service) ModelList(exclusions []string) ([]ModelInfo, error) {
	files, err := filepath.Glob(filepath.Join(GenBaseDir["model"], "*.go"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)

	excluded := make(map[string]struct{}, len(exclusions))
	for _, name := range exclusions {
		excluded[name] = struct{}{}
	}

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
				// 排除指定模型名
				if _, ok := excluded[typeSpec.Name.Name]; ok {
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
	sort.Slice(models, func(i, j int) bool { return models[i].Name < models[j].Name })
	return models, nil
}

// menuNode 菜单权限节点
type menuNode struct {
	title  string
	action string
}

// CreateMenuRule 写入后台菜单与权限节点，多层级路由如 a/b/c 会生成 a、a/b 两个目录与 a/b/c 菜单
func CreateMenuRule(ctx context.Context, repoAuth *repoAuth.AdminRuleRepository, basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) error {
	// 菜单 name/path 取前端路由路径，组件路径取视图目录
	views := basic["views"]
	menuPath := table.RoutePath
	if menuPath == "" {
		return nil
	}

	empty := ""
	keepalive := uint8(1)
	status := uint8(1)
	weigh := 0

	// 组装权限节点，作为菜单的子级
	nodes := []menuNode{
		{"查看", "read"},
		{"添加", "create"},
		{"更新", "update"},
		{"删除", "delete"},
	}
	if HasWeighDesign(fields) {
		nodes = append(nodes, menuNode{"排序", "sort"})
	}
	children := make([]*model.AdminRule, 0, len(nodes))
	for _, n := range nodes {
		children = append(children, &model.AdminRule{
			Type:      "node",
			Title:     n.title,
			Name:      menuPath + "/" + n.action,
			Path:      &empty,
			Icon:      &empty,
			OpenType:  &empty,
			URL:       &empty,
			Component: &empty,
			Keepalive: &keepalive,
			Extend:    &empty,
			Remark:    &empty,
			Weigh:     &weigh,
			Status:    &status,
		})
	}

	// 自底向上组装目录链: 中间层 Type=dir，最后一层 Type=menu，权限节点挂在菜单下
	component := "/" + strings.TrimPrefix(views.Dir+"/index.vue", "web/")
	segs := SplitPath(menuPath)
	for i := len(segs) - 1; i >= 0; i-- {
		levelPath := strings.Join(segs[:i+1], "/")
		isLast := i == len(segs)-1
		ruleType := "dir"
		title := segs[i]
		ruleComponent := &empty
		ruleOpenType := &empty
		if isLast {
			ruleType = "menu"
			title = table.Comment + "管理"
			ruleComponent = &component
			ruleOpenType = util.ToPtr("tab")
		}
		rule := &model.AdminRule{
			Type:      ruleType,
			Title:     title,
			Name:      levelPath,
			Path:      &levelPath,
			Icon:      &empty,
			OpenType:  ruleOpenType,
			URL:       &empty,
			Component: ruleComponent,
			Keepalive: &keepalive,
			Extend:    &empty,
			Remark:    &empty,
			Weigh:     &weigh,
			Status:    &status,
			Children:  children,
		}
		children = []*model.AdminRule{rule}
	}

	return repoAuth.BatchCreate(ctx, children, nil)
}

// HandleTableDesign 同步或创建数据表；返回执行的 SQL 列表
func HandleTableDesign(ctx context.Context, table dto.CRUDTable, fields []dto.CRUDFields) ([]string, error) {
	tableName := WithPrefix(table.Name)
	var sqls []string

	// 表不存在则按字段定义新建
	if !tbl.TableExists(ctx, tableName) {
		if err := tbl.CreateTable(ctx, tableName, table, fields, &sqls); err != nil {
			return nil, err
		}
		return sqls, nil
	}
	if err := tbl.SyncTableDesign(ctx, tableName, table, fields, &sqls); err != nil {
		return nil, err
	}
	return sqls, nil
}

// ParseFieldData 解析指定数据表的字段数据
func (s *Service) ParseFieldData(ctx context.Context, table string) ([]tbl.FieldItem, error) {
	return tbl.ParseFieldData(ctx, WithPrefix(table))
}
