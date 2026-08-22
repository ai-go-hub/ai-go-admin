package structx

import (
	"reflect"
	"strings"
	"sync"
)

// Field 结构体字段元数据
type Field struct {
	Index int               // 字段索引
	Key   string            // json tag 键名
	Type  reflect.Type      // 字段类型
	Tag   reflect.StructTag // 完整 struct tag
}

// fieldsCache 字段元数据切片缓存
// map[reflect.Type][]Field
var fieldsCache sync.Map

// FieldsOf 解析结构体的字段元数据
func FieldsOf(t reflect.Type) []Field {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if cached, ok := fieldsCache.Load(t); ok {
		return cached.([]Field)
	}

	var fields []Field
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			if sf.PkgPath != "" {
				// 未导出
				continue
			}
			fields = append(fields, Field{
				Index: i,
				Key:   jsonKey(sf),
				Type:  sf.Type,
				Tag:   sf.Tag,
			})
		}
	}
	actual, _ := fieldsCache.LoadOrStore(t, fields)
	return actual.([]Field)
}

// jsonKey 提取字段 json tag 的键名（忽略 ",options"），无 json tag 或 "-" 返回空串
func jsonKey(sf reflect.StructField) string {
	tag := sf.Tag.Get("json")
	if tag == "" || tag == "-" {
		return ""
	}
	if i := strings.Index(tag, ","); i >= 0 {
		tag = tag[:i]
	}
	return tag
}
