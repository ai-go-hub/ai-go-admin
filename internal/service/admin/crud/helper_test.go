package crud

import (
	"reflect"
	"testing"
)

func TestGenerateFileBasicData(t *testing.T) {
	cases := []struct {
		name  string
		typ   string
		table string
		app   string
		want  GenerateFileBasicDataInfo
	}{
		{
			"model 平铺",
			"model", "user_log", "admin",
			GenerateFileBasicDataInfo{Type: "model", Table: "user_log", App: "admin", Dir: "internal/model", File: "internal/model/user_log.go", Package: "model", LastName: "user_log", Name: "UserLog"},
		},
		{
			"handler",
			"handler", "auth_rule", "admin",
			GenerateFileBasicDataInfo{Type: "handler", Table: "auth_rule", App: "admin", Dir: "internal/handler/admin/auth", File: "internal/handler/admin/auth/rule.go", Package: "auth", LastName: "rule", Name: "AuthRule"},
		},
		{
			"service",
			"service", "auth_rule", "admin",
			GenerateFileBasicDataInfo{Type: "service", Table: "auth_rule", App: "admin", Dir: "internal/service/admin/auth", File: "internal/service/admin/auth/rule.go", Package: "auth", LastName: "rule", Name: "AuthRule"},
		},
		{
			"repository",
			"repository", "auth_rule", "admin",
			GenerateFileBasicDataInfo{Type: "repository", Table: "auth_rule", App: "admin", Dir: "internal/repository/admin/auth", File: "internal/repository/admin/auth/rule.go", Package: "auth", LastName: "rule", Name: "AuthRule"},
		},
		{
			"router",
			"router", "auth_rule", "admin",
			GenerateFileBasicDataInfo{Type: "router", Table: "auth_rule", App: "admin", Dir: "internal/router/admin/auth", File: "internal/router/admin/auth/rule.go", Package: "auth", LastName: "rule", Name: "AuthRule"},
		},
		{
			"views 目录",
			"views", "auth_rule", "admin",
			GenerateFileBasicDataInfo{Type: "views", Table: "auth_rule", App: "admin", Dir: "web/src/views/admin/auth/rule", LastName: "rule", Name: "AuthRule"},
		},
		{
			"lang 语言包",
			"lang", "auth_rule", "admin",
			GenerateFileBasicDataInfo{Type: "lang", Table: "auth_rule", App: "admin", Dir: "web/src/lang", LastName: "rule", Name: "AuthRule", CnFile: "web/src/lang/zh-cn/auth/rule.yaml", EnFile: "web/src/lang/en/auth/rule.yaml"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := GenerateFileBasicData(c.typ, c.table, c.app); got != c.want {
				t.Errorf("GenerateFileBasicData(%q, %q, %q) = %+v, want %+v", c.typ, c.table, c.app, got, c.want)
			}
		})
	}
}

func TestParseGenerateFileBasicData(t *testing.T) {
	t.Run("Generate 与 Parse 互逆", func(t *testing.T) {
		cases := []struct {
			typ, table, app, path string
		}{
			{"model", "auth_rule", "", "internal/model/auth_rule.go"},
			{"handler", "auth_rule", "admin", "internal/handler/admin/auth/rule.go"},
			{"service", "auth_rule", "admin", "internal/service/admin/auth/rule.go"},
			{"repository", "auth_rule", "admin", "internal/repository/admin/auth/rule.go"},
			{"router", "auth_rule", "admin", "internal/router/admin/auth/rule.go"},
			{"views", "auth_rule", "admin", "web/src/views/admin/auth/rule"},
		}
		for _, c := range cases {
			gen := GenerateFileBasicData(c.typ, c.table, c.app)
			if got := ParseGenerateFileBasicData(c.typ, c.path); got != gen {
				t.Errorf("%s: got %+v, want %+v", c.typ, got, gen)
			}
		}
	})

	checks := []struct {
		name string
		typ  string
		path string
		want GenerateFileBasicDataInfo
	}{
		{
			"model 平铺", "model", "internal/model/user_log.go",
			GenerateFileBasicDataInfo{Type: "model", Table: "user_log", Dir: "internal/model", File: "internal/model/user_log.go", Package: "model", LastName: "user_log", Name: "UserLog"},
		},
		{
			"model 子级", "model", "internal/model/auth/rule.go",
			GenerateFileBasicDataInfo{Type: "model", Table: "auth_rule", Dir: "internal/model/auth", File: "internal/model/auth/rule.go", Package: "auth", LastName: "auth_rule", Name: "AuthRule"},
		},
		{
			"lang 目录 zh-cn", "lang", "web/src/lang/zh-cn/auth",
			GenerateFileBasicDataInfo{Type: "lang", Table: "auth", Dir: "web/src/lang/zh-cn/auth", LastName: "auth", Name: "Auth"},
		},
		{
			"lang 目录 en", "lang", "web/src/lang/en/auth/rule",
			GenerateFileBasicDataInfo{Type: "lang", Table: "auth_rule", Dir: "web/src/lang/en/auth/rule", LastName: "rule", Name: "AuthRule"},
		},
		{
			"views 目录", "views", "web/src/views/admin/auth/rule",
			GenerateFileBasicDataInfo{Type: "views", Table: "auth_rule", App: "admin", Dir: "web/src/views/admin/auth/rule", LastName: "rule", Name: "AuthRule"},
		},
	}
	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			if got := ParseGenerateFileBasicData(c.typ, c.path); got != c.want {
				t.Errorf("ParseGenerateFileBasicData(%q, %q) = %+v, want %+v", c.typ, c.path, got, c.want)
			}
		})
	}

	t.Run("形态不符或未知类型返回零值", func(t *testing.T) {
		zeroCases := []struct{ typ, path string }{
			{"unknown", "anything"},                       // 未知类型
			{"handler", "C:/other/admin/auth/rule.go"},    // 前缀不匹配
			{"views", "web/src/views/admin/auth/rule.go"}, // views 传 .go 文件
			{"lang", "web/src/lang/zh-cn/auth/rule.yaml"}, // lang 传 .yaml 文件
			{"handler", "internal/handler/admin/auth"},    // handler 传目录
			{"model", "internal/model/auth"},              // model 传目录
			{"model", "internal/model"},                   // base 目录非 .go
		}
		for _, c := range zeroCases {
			if got := ParseGenerateFileBasicData(c.typ, c.path); got != (GenerateFileBasicDataInfo{}) {
				t.Errorf("%s/%s: got %+v, want 零值", c.typ, c.path, got)
			}
		}
	})
}

func TestSplitTablePath(t *testing.T) {
	cases := []struct {
		name  string
		table string
		want  []string
	}{
		{"下划线多段", "user_log", []string{"user", "log"}},
		{"下划线多段2", "auth_rule", []string{"auth", "rule"}},
		{"斜杠多段", "auth/rule", []string{"auth", "rule"}},
		{"空串", "", []string{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := SplitTablePath(c.table); !reflect.DeepEqual(got, c.want) {
				t.Errorf("SplitTablePath(%q) = %v, want %v", c.table, got, c.want)
			}
		})
	}
}

func TestPascalCase(t *testing.T) {
	cases := []struct {
		name string
		segs []string
		want string
	}{
		{"user_log 拆分后", []string{"user", "log"}, "UserLog"},
		{"auth_rule 拆分后", []string{"auth", "rule"}, "AuthRule"},
		{"空段被跳过", []string{"user", "", "log"}, "UserLog"},
		{"空输入", nil, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := PascalCase(c.segs); got != c.want {
				t.Errorf("PascalCase(%v) = %q, want %q", c.segs, got, c.want)
			}
		})
	}
}

func TestBuildParts(t *testing.T) {
	cases := []struct {
		name  string
		parts []string
		want  []string
	}{
		{"过滤空段", []string{"a", "", "b"}, []string{"a", "b"}},
		{"全空段", []string{"", ""}, []string{}},
		{"nil", nil, []string{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := BuildParts(c.parts); !reflect.DeepEqual(got, c.want) {
				t.Errorf("BuildParts(%v) = %v, want %v", c.parts, got, c.want)
			}
		})
	}
}

func TestLastDirName(t *testing.T) {
	cases := []struct {
		name string
		path string
		want string
	}{
		{"两级", "internal/model", "model"},
		{"多级", "web/src/views/admin/auth", "auth"},
		{"单段", "auth", "auth"},
		{"空串", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := LastDirName(c.path); got != c.want {
				t.Errorf("LastDirName(%q) = %q, want %q", c.path, got, c.want)
			}
		})
	}
}
