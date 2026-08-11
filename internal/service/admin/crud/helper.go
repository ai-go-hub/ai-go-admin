package crud

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"
	tbl "github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud/table"
	"github.com/ai-go-hub/ai-go-admin/pkg/filesystem"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"
)

// GenBaseDir 各类生成文件的基准目录
var GenBaseDir = map[string]string{
	"model":      `internal/model`,
	"handler":    `internal/handler`,
	"service":    `internal/service`,
	"repository": `internal/repository`,
	"router":     `internal/router`,
	"views":      `web/src/views`,
	"lang":       `web/src/lang`,
}

// ModulePath 项目模块路径
const ModulePath = "github.com/ai-go-hub/ai-go-admin"

// GenerateRoutePath 获取生成路由路径
func GenerateRoutePath(table string) string {
	segs := SplitPath(table)
	if len(segs) == 0 {
		return ""
	}
	return strings.Join(segs, "/")
}

// GenerateFileBasicData 获取生成文件的基本信息
func GenerateFileBasicData(typ, path, app string) dto.GenerateFileBasicDataInfo {
	data := dto.GenerateFileBasicDataInfo{
		Type: typ,
		Path: path,
		App:  app,
	}

	segs := SplitPath(path)
	if len(segs) == 0 {
		return data
	}
	data.LastName = segs[len(segs)-1]
	data.Name = util.SnakeToPascal(path)

	switch typ {
	case "model":
		// 模型平铺在 internal/model 目录，和迁移保持一致
		data.Dir = GenBaseDir["model"]
		data.LastName = strings.Join(segs, "_")
		data.File = data.Dir + "/" + data.LastName + ".go"
		data.Package = "model"
	case "handler", "service", "repository", "router":
		dir := strings.Join(BuildParts(append([]string{GenBaseDir[typ], app}, segs[:len(segs)-1]...)), "/")
		data.Dir = dir
		data.File = dir + "/" + data.LastName + ".go"
		data.Package = LastDirName(dir)
	case "views":
		data.Dir = strings.Join(BuildParts(append([]string{GenBaseDir[typ], app}, segs...)), "/")
	case "lang":
		rel := strings.Join(BuildParts(segs), "/") + ".yaml"
		data.Dir = GenBaseDir["lang"]
		data.CnFile = data.Dir + "/zh-cn/" + rel
		data.EnFile = data.Dir + "/en/" + rel
		data.LangKey = strings.Join(BuildParts(segs), ".")
	}

	return data
}

// ParseGenerateFileBasicData 根据类型从 GenerateFileBasicData 生成的目录或文件路径反向解析基本信息
func ParseGenerateFileBasicData(typ, path string) dto.GenerateFileBasicDataInfo {
	data := dto.GenerateFileBasicDataInfo{Type: typ}

	base, ok := GenBaseDir[typ]
	if !ok || (path != base && !strings.HasPrefix(path, base+"/")) {
		return dto.GenerateFileBasicDataInfo{}
	}
	rest := strings.TrimPrefix(strings.TrimPrefix(path, base), "/")

	// 输入形态确定: views/lang 只支持目录输入，其余只支持 .go 文件输入
	if typ == "views" || typ == "lang" {
		if strings.HasSuffix(rest, ".go") || strings.HasSuffix(rest, ".yaml") {
			return dto.GenerateFileBasicDataInfo{}
		}
	} else if !strings.HasSuffix(rest, ".go") {
		return dto.GenerateFileBasicDataInfo{}
	}

	// 统一按 / 拆段，各分支只处理特有逻辑
	segs := BuildParts(strings.Split(rest, "/"))

	switch typ {
	case "views":
		data.Dir = path
		if len(segs) > 0 {
			data.App = segs[0]
			segs = segs[1:]
		}
	case "lang":
		data.Dir = path
		if len(segs) > 0 && (segs[0] == "zh-cn" || segs[0] == "en") {
			segs = segs[1:]
		} else {
			segs = nil
		}
	case "model":
		// model
		dirSegs := segs[:len(segs)-1]
		segs = append(dirSegs, SplitPath(strings.TrimSuffix(segs[len(segs)-1], ".go"))...)
		data.Dir = path[:strings.LastIndex(path, "/")]
		data.Package = LastDirName(data.Dir)
		data.File = path
	default:
		// handler/service/repository/router
		data.App = segs[0]
		segs = segs[1:]
		if len(segs) == 0 {
			return dto.GenerateFileBasicDataInfo{}
		}
		segs[len(segs)-1] = strings.TrimSuffix(segs[len(segs)-1], ".go")
		data.Dir = path[:strings.LastIndex(path, "/")]
		data.Package = LastDirName(data.Dir)
		data.File = path
	}

	// 统一设置表名信息
	if len(segs) > 0 {
		data.Path = strings.Join(segs, "_")
		if typ == "model" {
			data.LastName = data.Path
		} else {
			data.LastName = segs[len(segs)-1]
		}
		data.Name = util.SnakeToPascal(data.Path)
	}

	return data
}

// DictItem 字段字典项
type DictItem struct {
	Key   string // 字典 key，如 opt0
	Value string // 字典中文值，如 选项一
}

// ParseFieldComment 解析字段注释，格式: 标题[:key=值,key=值..]
func ParseFieldComment(comment string) (title string, dict []DictItem) {
	title, after, found := strings.Cut(comment, ":")
	if !found {
		return comment, nil
	}
	for item := range strings.SplitSeq(after, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		key, val, ok := strings.Cut(item, "=")
		if !ok {
			continue
		}
		dict = append(dict, DictItem{Key: strings.TrimSpace(key), Value: strings.TrimSpace(val)})
	}
	return title, dict
}

// FormValidators 返回字段 Form.validator 中的规则名列表（如 required、email）
func FormValidators(f dto.CRUDFields) []string {
	raw, ok := f.Form["validator"]
	if !ok {
		return nil
	}
	validators, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(validators))
	for _, v := range validators {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

// HasValidatorRule 判断字段 Form.validator 是否包含指定规则
func HasValidatorRule(f dto.CRUDFields, name string) bool {
	return slices.Contains(FormValidators(f), name)
}

// SplitPath 将路径按 _ 或 / 切分
func SplitPath(table string) []string {
	return strings.FieldsFunc(table, func(r rune) bool {
		return r == '_' || r == '/'
	})
}

// BuildParts 拼接路径段，过滤空段
func BuildParts(parts []string) []string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// LastDirName 取路径最后一个目录名
func LastDirName(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[i+1:]
	}
	return path
}

// WriteGeneratedFile 将内容写入项目内的生成文件（自动创建目录），Go 文件写入后自动执行 go fmt
func WriteGeneratedFile(relPath, content string) error {
	path, ok := filesystem.AbsInProject(relPath)
	if !ok {
		return errors.New("生成文件路径不合法: " + relPath)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return err
	}
	if strings.EqualFold(filepath.Ext(path), ".go") {
		return filesystem.FormatGoFile(path)
	}
	return nil
}

// RenderTmpl 执行模板并返回渲染结果
func RenderTmpl(t *template.Template, data any) (string, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ModuleImport 模块内导入信息
type ModuleImport struct {
	Alias string // 导入别名，空串表示无
	Path  string // 导入路径
	Pkg   string // 目录包名
}

// ModuleImportOf 组装模块内导入信息，
// prefix 为别名前缀（repo/svc/handler）
func ModuleImportOf(info dto.GenerateFileBasicDataInfo, prefix string) ModuleImport {
	return ModuleImport{
		Alias: prefix + util.SnakeToPascal(info.Package),
		Path:  ModulePath + "/" + info.Dir,
		Pkg:   info.Package,
	}
}

// ImportSpec 生成 import 行，alias 为可选别名，省略时为无别名导入
func ImportSpec(importPath string, alias ...string) string {
	if len(alias) == 0 {
		return `"` + importPath + `"`
	}
	return alias[0] + ` "` + importPath + `"`
}

// BuildImportBlock 组装 import 块
func BuildImportBlock(specs []string) string {
	if len(specs) == 0 {
		return ""
	}

	std, internal, third := make([]string, 0, len(specs)), make([]string, 0, len(specs)), make([]string, 0, len(specs))
	seen := make(map[string]struct{}, len(specs))
	for _, spec := range specs {
		start, end := strings.IndexByte(spec, '"'), strings.LastIndexByte(spec, '"')
		path := spec
		if start >= 0 && end > start {
			path = spec[start+1 : end]
		}
		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}
		switch {
		case strings.HasPrefix(path, ModulePath+"/"):
			internal = append(internal, spec)
		case strings.Contains(path, "."):
			third = append(third, spec)
		default:
			std = append(std, spec)
		}
	}

	// 组内顺序生成后由 go fmt 自动整理
	groups := [][]string{std, internal, third}
	var b strings.Builder
	b.WriteString("import (\n")
	first := true
	for _, group := range groups {
		if len(group) == 0 {
			continue
		}
		if !first {
			b.WriteString("\n")
		}
		first = false
		for _, s := range group {
			fmt.Fprintf(&b, "\t%s\n", s)
		}
	}
	b.WriteString(")")
	return b.String()
}

// WithPrefix 规范化表名，已带前缀则先剔除再拼接，防止前缀重复
func WithPrefix(name string) string {
	prefix := config.Get().Database.Prefix
	return prefix + strings.TrimPrefix(name, prefix)
}

// TableFieldComments 获取指定数据表 字段名 -> 注释 的映射
func TableFieldComments(ctx context.Context, table string) (map[string]string, error) {
	db := database.DB().WithContext(ctx)
	table = WithPrefix(table)

	var fields []FieldInfo
	if err := db.Raw(`
SELECT
    a.attname AS "name",
    COALESCE(col_description(a.attrelid, a.attnum), '') AS "comment"
FROM pg_catalog.pg_attribute a
JOIN pg_catalog.pg_class c ON c.oid = a.attrelid
JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = current_schema()
    AND c.relname = ?
    AND a.attnum > 0
    AND NOT a.attisdropped`, table).Scan(&fields).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string, len(fields))
	for _, f := range fields {
		result[f.Name] = f.Comment
	}
	return result, nil
}

// TableComment 获取指定数据表的注释
func TableComment(ctx context.Context, table string) (string, error) {
	db := database.DB().WithContext(ctx)
	table = WithPrefix(table)

	var comment string
	if err := db.Raw(`
SELECT COALESCE(obj_description(t.oid, 'pg_class'), '') AS "comment"
FROM pg_catalog.pg_class t
JOIN pg_catalog.pg_namespace n ON n.oid = t.relnamespace
WHERE n.nspname = current_schema()
    AND t.relname = ?`, table).Scan(&comment).Error; err != nil {
		return "", err
	}
	return comment, nil
}

// TableExists 检查数据表是否存在
func (s *Service) TableExists(ctx context.Context, table string) (bool, error) {
	db := database.DB().WithContext(ctx)
	table = WithPrefix(table)

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

// TableEmpty 判断数据表是否为空（无数据行）；表不存在时视为空
func (s *Service) TableEmpty(ctx context.Context, table string) (bool, error) {
	table = WithPrefix(table)
	if !tbl.TableExists(ctx, table) {
		return true, nil
	}
	var count int64
	if err := database.DB().WithContext(ctx).Raw("SELECT COUNT(*) FROM " + tbl.QuoteIdent(table)).Scan(&count).Error; err != nil {
		return true, err
	}
	return count == 0, nil
}

// QuoteStrSlice 生成 []string{"a", "b"} 字面量
func QuoteStrSlice(items []string) string {
	quoted := make([]string, len(items))
	for i, s := range items {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	return "[]string{" + strings.Join(quoted, ", ") + "}"
}

// PkFieldName 获取主键字段名称
func PkFieldName(fields []dto.CRUDFields) string {
	for _, f := range fields {
		if f.PrimaryKey {
			return f.Name
		}
	}
	return ""
}

// HasRemoteSelectsDesign 判断字段列表内是否有远程多选类型的设计
func HasRemoteSelectsDesign(fields []dto.CRUDFields) bool {
	for _, f := range fields {
		if f.DesignType == "remoteSelects" {
			return true
		}
	}
	return false
}

// HasWeighDesign 判断字段列表内是否有权重类型的设计
func HasWeighDesign(fields []dto.CRUDFields) bool {
	for _, f := range fields {
		if f.DesignType == "weigh" {
			return true
		}
	}
	return false
}
