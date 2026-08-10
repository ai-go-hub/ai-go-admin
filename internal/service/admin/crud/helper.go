package crud

import (
	"context"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/database"
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

// GenerateFileBasicDataInfo 生成文件的基本信息
type GenerateFileBasicDataInfo struct {
	Type     string `json:"type"`              // 生成文件类型: model/handler/service/repository/router/views/lang
	Table    string `json:"table"`             // 表名
	App      string `json:"app"`               // handler/service/repository/router/views 的一级子目录，如 admin、common
	Dir      string `json:"dir"`               // 生成目录
	File     string `json:"file,omitempty"`    // 文件完整路径
	Package  string `json:"package,omitempty"` // go 文件 package 值
	LastName string `json:"last_name"`         // 表名最后一级，如 user_log 的 log
	Name     string `json:"name"`              // go 模型名，表名转大驼峰，如 user_log 的 UserLog
	CnFile   string `json:"cn_file,omitempty"` // 语言包 zh-cn 文件路径
	EnFile   string `json:"en_file,omitempty"` // 语言包 en 文件路径
}

// GenerateFileBasicData 获取生成文件的基本信息
func GenerateFileBasicData(typ, table, app string) GenerateFileBasicDataInfo {
	data := GenerateFileBasicDataInfo{
		Type:  typ,
		Table: table,
		App:   app,
	}

	segs := SplitTablePath(table)
	if len(segs) == 0 {
		return data
	}
	data.LastName = segs[len(segs)-1]
	data.Name = PascalCase(segs)

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
		rel := strings.Join(segs, "/") + ".yaml"
		data.Dir = GenBaseDir["lang"]
		data.CnFile = data.Dir + "/zh-cn/" + rel
		data.EnFile = data.Dir + "/en/" + rel
	}

	return data
}

// ParseGenerateFileBasicData 根据类型从 GenerateFileBasicData 生成的目录或文件路径反向解析基本信息
func ParseGenerateFileBasicData(typ, path string) GenerateFileBasicDataInfo {
	data := GenerateFileBasicDataInfo{Type: typ}

	base, ok := GenBaseDir[typ]
	if !ok || (path != base && !strings.HasPrefix(path, base+"/")) {
		return GenerateFileBasicDataInfo{}
	}
	rest := strings.TrimPrefix(strings.TrimPrefix(path, base), "/")

	// 输入形态确定: views/lang 只支持目录输入，其余只支持 .go 文件输入
	if typ == "views" || typ == "lang" {
		if strings.HasSuffix(rest, ".go") || strings.HasSuffix(rest, ".yaml") {
			return GenerateFileBasicDataInfo{}
		}
	} else if !strings.HasSuffix(rest, ".go") {
		return GenerateFileBasicDataInfo{}
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
		segs = append(dirSegs, SplitTablePath(strings.TrimSuffix(segs[len(segs)-1], ".go"))...)
		data.Dir = path[:strings.LastIndex(path, "/")]
		data.Package = LastDirName(data.Dir)
		data.File = path
	default:
		// handler/service/repository/router
		data.App = segs[0]
		segs = segs[1:]
		if len(segs) == 0 {
			return GenerateFileBasicDataInfo{}
		}
		segs[len(segs)-1] = strings.TrimSuffix(segs[len(segs)-1], ".go")
		data.Dir = path[:strings.LastIndex(path, "/")]
		data.Package = LastDirName(data.Dir)
		data.File = path
	}

	// 统一设置表名信息
	if len(segs) > 0 {
		data.Table = strings.Join(segs, "_")
		if typ == "model" {
			data.LastName = data.Table
		} else {
			data.LastName = segs[len(segs)-1]
		}
		data.Name = PascalCase(segs)
	}

	return data
}

// SplitTablePath 将表名按 _ 或 / 切割
func SplitTablePath(table string) []string {
	return strings.FieldsFunc(table, func(r rune) bool {
		return r == '_' || r == '/'
	})
}

// PascalCase 将切割后的表名转为大驼峰，如 user_log -> UserLog
func PascalCase(segs []string) string {
	var b strings.Builder
	for _, s := range segs {
		if s == "" {
			continue
		}
		b.WriteString(strings.ToUpper(s[:1]))
		b.WriteString(s[1:])
	}
	return b.String()
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

// TableComment 获取指定数据表的注释
func TableComment(ctx context.Context, table string) (string, error) {
	db := database.DB().WithContext(ctx)
	table = withPrefix(table)

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
