package gen

import (
	"strings"
	"testing"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"
)

// modelBasic 组装模型文件基础数据（测试用，仅需模型名）
func modelBasic(name string) map[string]dto.GenerateFileBasicDataInfo {
	return map[string]dto.GenerateFileBasicDataInfo{
		"model": {Name: name},
	}
}

// TestRenderModelFile 验证模型文件数据组装与模板渲染结果
func TestRenderModelFile(t *testing.T) {
	table := dto.CRUDTable{
		Name:    "user_log",
		Comment: "用户日志",
	}
	fields := []dto.CRUDFields{
		{Name: "id", Type: "bigint", Comment: "ID", PrimaryKey: true, Generated: "GENERATED ALWAYS"},
		{Name: "name", Type: "varchar", Length: 255, Comment: "名称", Null: true},
		{Name: "status", Type: "smallint", Comment: "状态", Null: false, Default: "1", DefaultType: "INPUT"},
		{Name: "tags", Type: "jsonb", Comment: "标签", Null: true},
		{Name: "created_at", Type: "timestamptz", Comment: "创建时间", Null: true},
	}

	data := buildModelFileData(modelBasic("UserLog"), table, fields)

	if data.Name != "UserLog" {
		t.Errorf("Name = %q, want UserLog", data.Name)
	}
	if data.Table != "user_log" {
		t.Errorf("Table = %q, want user_log", data.Table)
	}

	// 主键不加指针，可空字段用指针
	typeCases := []struct{ name, want string }{
		{"ID", "uint"},
		{"Name", "*string"},
		{"Status", "int16"},
		{"Tags", "*datatypes.JSON"},
		{"CreatedAt", "*time.Time"},
	}
	if len(data.Fields) != len(typeCases) {
		t.Fatalf("字段数量 = %d, want %d", len(data.Fields), len(typeCases))
	}
	for i, tc := range typeCases {
		if data.Fields[i].Name != tc.name || data.Fields[i].Type != tc.want {
			t.Errorf("字段[%d] = %s %s, want %s %s", i, data.Fields[i].Name, data.Fields[i].Type, tc.name, tc.want)
		}
	}

	// import 块应包含 time 与 gorm.io/datatypes 且按 std/第三方 分组
	if !strings.Contains(data.Imports, `"time"`) || !strings.Contains(data.Imports, `"gorm.io/datatypes"`) {
		t.Errorf("Imports 不完整: %q", data.Imports)
	}

	content, err := crud.RenderTmpl(modelTmpl, data)
	if err != nil {
		t.Fatal(err)
	}
	content = strings.ReplaceAll(content, "\r\n", "\n")

	for _, want := range []string{
		"package model\n",
		"import (\n\t\"time\"\n\n\t\"gorm.io/datatypes\"\n)\n",
		"// UserLog 用户日志\n",
		"type UserLog struct {\n",
		"`gorm:\"comment:ID;type:bigint;primarykey;autoIncrement\" json:\"id\"`",
		"`gorm:\"comment:名称;type:varchar(255)\" json:\"name\"`",
		"`gorm:\"comment:状态;type:smallint;not null;default:1\" json:\"status\"`",
		"`gorm:\"comment:标签;type:jsonb\" json:\"tags\"`",
		"`gorm:\"comment:创建时间;type:timestamptz\" json:\"created_at\"`",
		"func (UserLog) TableName() string {\n\treturn \"user_log\"\n}\n",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("渲染结果缺少 %q\n---\n%s", want, content)
		}
	}
	if !strings.HasSuffix(content, "\n") {
		t.Error("渲染结果未以换行结尾")
	}
	t.Logf("渲染结果:\n%s", content)
}

// TestRenderModelFileNoImports 验证无 time/datatypes 字段时不渲染 import 块
func TestRenderModelFileNoImports(t *testing.T) {
	table := dto.CRUDTable{Name: "test2", Comment: "测试"}
	fields := []dto.CRUDFields{
		{Name: "id", Type: "bigint", Comment: "ID", PrimaryKey: true},
		{Name: "name", Type: "varchar", Length: 64, Comment: "名称", Null: true},
	}
	content, err := crud.RenderTmpl(modelTmpl, buildModelFileData(modelBasic("Test2"), table, fields))
	if err != nil {
		t.Fatal(err)
	}
	content = strings.ReplaceAll(content, "\r\n", "\n")
	if strings.Contains(content, "import") {
		t.Errorf("无 import 时不应出现 import 块:\n%s", content)
	}
	for _, want := range []string{
		"package model\n\n// Test2 测试\n",
		"type Test2 struct {\n",
		"func (Test2) TableName() string {\n\treturn \"test2\"\n}\n",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("渲染结果缺少 %q\n---\n%s", want, content)
		}
	}
}

// TestBuildGoType 验证字段 Go 类型构建
func TestBuildGoType(t *testing.T) {
	cases := []struct {
		name  string
		field dto.CRUDFields
		want  string
	}{
		{"主键整型", dto.CRUDFields{PrimaryKey: true, Type: "bigint"}, "uint"},
		{"主键字符", dto.CRUDFields{PrimaryKey: true, Type: "varchar"}, "string"},
		{"varchar 必填", dto.CRUDFields{Type: "varchar"}, "string"},
		{"varchar 可空", dto.CRUDFields{Type: "varchar", Null: true}, "*string"},
		{"bigint 必填", dto.CRUDFields{Type: "bigint"}, "int64"},
		{"bigint 可空", dto.CRUDFields{Type: "bigint", Null: true}, "*int64"},
		{"smallint 必填", dto.CRUDFields{Type: "smallint"}, "int16"},
		{"smallint 可空", dto.CRUDFields{Type: "smallint", Null: true}, "*int16"},
		{"boolean 必填", dto.CRUDFields{Type: "boolean"}, "bool"},
		{"boolean 可空", dto.CRUDFields{Type: "boolean", Null: true}, "*bool"},
		{"numeric 必填", dto.CRUDFields{Type: "numeric"}, "float64"},
		{"numeric 可空", dto.CRUDFields{Type: "numeric", Null: true}, "*float64"},
		{"jsonb 必填", dto.CRUDFields{Type: "jsonb"}, "*datatypes.JSON"},
		{"jsonb 可空", dto.CRUDFields{Type: "jsonb", Null: true}, "*datatypes.JSON"},
		{"timestamptz 必填", dto.CRUDFields{Type: "timestamptz"}, "time.Time"},
		{"timestamptz 可空", dto.CRUDFields{Type: "timestamptz", Null: true}, "*time.Time"},
		{"year 必填", dto.CRUDFields{Type: "smallint", DesignType: "year"}, "string"},
		{"year 可空", dto.CRUDFields{Type: "smallint", DesignType: "year", Null: true}, "*string"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := buildGoType(c.field); got != c.want {
				t.Errorf("buildGoType(%+v) = %q, want %q", c.field, got, c.want)
			}
		})
	}
}

// TestRenderModelFileXssTag 验证 editor 字段自动追加 xss:"html" 标签，非 editor 字段不加
func TestRenderModelFileXssTag(t *testing.T) {
	table := dto.CRUDTable{Name: "article", Comment: "文章"}
	fields := []dto.CRUDFields{
		{Name: "id", Type: "bigint", Comment: "ID", PrimaryKey: true},
		{Name: "title", Type: "varchar", Length: 255, Comment: "标题"},
		{Name: "content", Type: "text", Comment: "内容", DesignType: "editor"},
	}
	content, err := crud.RenderTmpl(modelTmpl, buildModelFileData(modelBasic("Article"), table, fields))
	if err != nil {
		t.Fatal(err)
	}
	content = strings.ReplaceAll(content, "\r\n", "\n")
	if !strings.Contains(content, `json:"content" xss:"html"`) {
		t.Errorf("editor 字段缺少 xss:\"html\" 标签\n%s", content)
	}
	if strings.Contains(content, `json:"title" xss:`) {
		t.Errorf("非 editor 字段不应有 xss tag\n%s", content)
	}
}
