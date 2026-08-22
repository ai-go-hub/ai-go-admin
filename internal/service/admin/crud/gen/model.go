package gen

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"
	"github.com/ai-go-hub/ai-go-admin/pkg/filesystem"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"
)

//go:embed tmpl/model.tmpl
var modelTmplStr string

// modelTmpl 模型文件内容模板
var modelTmpl = template.Must(template.New("model").Parse(modelTmplStr))

// modelFileData 模型文件模板渲染数据
type modelFileData struct {
	Name      string           // 模型名
	Table     string           // 表名
	Comment   string           // 表注释
	Imports   string           // import 块内容，空串表示无 import
	Fields    []modelFileField // 字段行
	HasWeigh  bool             // 是否含权重字段，含则生成 AfterCreate 回填权重
	WeighType string           // 权重字段 Go 基础类型（不含指针）
}

// modelFileField 模型字段行
type modelFileField struct {
	Name       string // Go 字段名
	Type       string // Go 类型
	JSONTag    string // json tag 值
	GormTag    string // gorm tag 值
	BindingTag string // binding tag 值
	XssTag     string // xss tag 值
	Line       string // 单行字段定义
}

// buildModelFileData 组装模型文件模板数据
func buildModelFileData(basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) modelFileData {
	data := modelFileData{
		Name:    basic["model"].Name,
		Table:   table.Name,
		Comment: table.Comment,
	}
	modelPkg := basic["model"].Package

	imports := make([]string, 0, 2)
	fileFields := make([]modelFileField, 0, len(fields))
	for _, f := range fields {
		switch goTypeByDBType(f.Type) {
		case "time.Time":
			imports = append(imports, crud.ImportSpec("time"))
		case "datatypes.JSON":
			imports = append(imports, crud.ImportSpec("gorm.io/datatypes"))
		}

		mf := modelFileField{
			Name:       snakeToPascal(f.Name),
			Type:       buildGoType(f),
			JSONTag:    f.Name,
			GormTag:    buildGormTag(f),
			BindingTag: buildBindingTag(f),
		}
		if f.DesignType == "editor" {
			mf.XssTag = "html"
		}

		fileFields = append(fileFields, mf)

		if f.DesignType == "weigh" {
			data.HasWeigh = true
			data.WeighType = goTypeByDBType(f.Type)
		}

		// remoteSelect 字段向模型定义追加 belongs to 关联字段
		if f.DesignType == "remoteSelect" {
			if assoc, ok := buildModelAssociationField(f, modelPkg); ok {
				fileFields = append(fileFields, assoc)
				if info := parseRemoteInfo(f); info.ModelPkg != "" && info.ModelPkg != modelPkg && info.ModelFile != "" {
					imports = append(imports, crud.ImportSpec(crud.ModulePath+"/"+filesystem.Dir(info.ModelFile)))
				}
			}
		}

		// remoteSelects 字段向模型定义追加只读展示字段
		if f.DesignType == "remoteSelects" {
			if info := parseRemoteInfo(f); info.ModelName != "" && info.Field != "" {
				name := RemoteSelectsListFieldName(f)
				fileFields = append(fileFields, modelFileField{
					Name:    name,
					Type:    "[]string",
					JSONTag: util.PascalToSnake(name),
					GormTag: "-",
				})
			}
		}
	}
	if data.HasWeigh {
		imports = append(imports, crud.ImportSpec("gorm.io/gorm"))
	}
	data.Imports = crud.BuildImportBlock(imports)

	for i := range fileFields {
		f := &fileFields[i]
		line := fmt.Sprintf("\t%s %s `gorm:\"%s\" json:\"%s\"", f.Name, f.Type, f.GormTag, f.JSONTag)
		if f.BindingTag != "" {
			line += fmt.Sprintf(" binding:\"%s\"", f.BindingTag)
		}
		if f.XssTag != "" {
			line += fmt.Sprintf(" xss:\"%s\"", f.XssTag)
		}
		f.Line = line + "`"
	}
	data.Fields = fileFields
	return data
}

// snakeToPascal 模型生成专用的特殊版 snake_case 转大驼峰，
// 如: id -> ID,updated_at -> UpdatedAt
func snakeToPascal(name string) string {
	segs := strings.FieldsFunc(name, func(r rune) bool { return r == '_' })
	var b strings.Builder
	for _, seg := range segs {
		b.WriteString(fieldSegName(seg))
	}
	return b.String()
}

// fieldSegName 单个字段名段转 Go 名称，常见缩写整体大写
func fieldSegName(seg string) string {
	switch seg {
	case "id", "uuid", "url", "ip", "api":
		return strings.ToUpper(seg)
	}
	if seg == "" {
		return seg
	}
	return strings.ToUpper(seg[:1]) + seg[1:]
}

// goTypeByDBType 数据库字段类型对应的 Go 基础类型（不含指针）
func goTypeByDBType(dbType string) string {
	switch dbType {
	case "varchar", "char", "text", "uuid":
		return "string"
	case "bigint":
		return "int64"
	case "integer":
		return "int"
	case "smallint":
		return "int16"
	case "boolean":
		return "bool"
	case "numeric", "real", "double precision":
		return "float64"
	case "json", "jsonb":
		return "datatypes.JSON"
	case "date", "timestamp", "timestamptz", "time", "timetz":
		return "time.Time"
	default:
		return "string"
	}
}

// buildGoType 构建字段的 Go 类型
func buildGoType(f dto.CRUDFields) string {
	if f.PrimaryKey {
		// 主键为整型时用 uint 且不加指针
		if isIntegerType(f.Type) {
			return "uint"
		}
		return goTypeByDBType(f.Type)
	}
	// year 接受为 *string 但数据库内存储为 smallint
	base := goTypeByDBType(f.Type)
	if f.DesignType == "year" {
		base = "string"
	}
	// jsonb 值类型会把 null 存成字符串 "null"，统一使用指针强行规避错误设定
	if base == "datatypes.JSON" {
		return "*" + base
	}
	if f.Null {
		return "*" + base
	}
	return base
}

// isIntegerType 判断是否为整型字段类型
func isIntegerType(dbType string) bool {
	switch dbType {
	case "bigint", "integer", "smallint":
		return true
	}
	return false
}

// buildGormTag 构建字段的 gorm tag
func buildGormTag(f dto.CRUDFields) string {
	parts := make([]string, 0, 6)
	if f.Comment != "" {
		parts = append(parts, "comment:"+f.Comment)
	}
	if typ := buildDBType(f); typ != "" {
		parts = append(parts, "type:"+typ)
	}
	if f.PrimaryKey {
		parts = append(parts, "primarykey")
		if isIntegerType(f.Type) {
			parts = append(parts, "autoIncrement")
		}
	}
	if f.Unique {
		parts = append(parts, "uniqueIndex")
	}
	if !f.PrimaryKey && !f.Null {
		parts = append(parts, "not null")
	}
	if def := buildGormDefault(f); def != "" {
		parts = append(parts, "default:"+def)
	}
	return strings.Join(parts, ";")
}

// buildBindingTag 构建字段的 binding tag
func buildBindingTag(f dto.CRUDFields) string {
	// 当字段 Form.validator 含 required 时返回 required，用于模型 binding tag
	if crud.HasValidatorRule(f, "required") {
		return "required"
	}
	return ""
}

// buildDBType 构建字段的 type 定义
func buildDBType(f dto.CRUDFields) string {
	switch f.Type {
	case "varchar", "char":
		if f.Length > 0 {
			return fmt.Sprintf("%s(%d)", f.Type, f.Length)
		}
		return f.Type
	case "numeric":
		if f.Length > 0 && f.Precision > 0 {
			return fmt.Sprintf("numeric(%d,%d)", f.Length, f.Precision)
		}
		return f.Type
	default:
		return f.Type
	}
}

// buildGormDefault 构建字段的 default 值
func buildGormDefault(f dto.CRUDFields) string {
	switch f.DefaultType {
	case "EMPTY STRING":
		return "''"
	case "INPUT":
		if goTypeByDBType(f.Type) == "string" {
			return "'" + f.Default + "'"
		}
		return f.Default
	}
	return ""
}

// CreateModelFile 创建模型文件
func CreateModelFile(basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) error {
	data := buildModelFileData(basic, table, fields)
	content, err := crud.RenderTmpl(modelTmpl, data)
	if err != nil {
		return err
	}
	return crud.WriteGeneratedFile(basic["model"].File, content)
}
