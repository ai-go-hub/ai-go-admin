package structx

import (
	"reflect"
	"testing"
)

// metaModel 测试用模型: 含普通字段、json:"-"、无 json tag、带 xss 扩展 tag
type metaModel struct {
	ID      uint    `json:"id"`
	Name    *string `json:"name"`
	Secret  string  `json:"-"`
	Content string  `json:"content" xss:"html"`
	Raw     string
}

func TestFieldsOf(t *testing.T) {
	fields := FieldsOf(reflect.TypeFor[metaModel]())

	keys := make(map[string]Field)
	for _, f := range fields {
		keys[f.Key] = f
	}
	for _, key := range []string{"id", "name", "content"} {
		f, ok := keys[key]
		if !ok {
			t.Errorf("缺少字段 %q", key)
			continue
		}
		if f.Index < 0 {
			t.Errorf("字段 %q Index 无效: %d", key, f.Index)
		}
	}
	// 扩展 tag 可读取
	if f := keys["content"]; f.Tag.Get("xss") != "html" {
		t.Errorf("content 的 xss tag = %q, want html", f.Tag.Get("xss"))
	}
	// json:"-" 与无 json tag 的字段 Key 应为空
	for _, f := range fields {
		if f.Tag.Get("json") == "-" && f.Key != "" {
			t.Errorf("json:\"-\" 字段 Key 应为空, got %q", f.Key)
		}
		if f.Tag.Get("json") == "" && f.Key != "" {
			t.Errorf("无 json tag 字段 Key 应为空, got %q", f.Key)
		}
	}
}

func TestFieldsOfPointerAndNonStruct(t *testing.T) {
	if got := FieldsOf(reflect.TypeFor[*metaModel]()); len(got) == 0 {
		t.Error("指针类型应解引用后解析")
	}
	if got := FieldsOf(reflect.TypeOf(123)); len(got) != 0 {
		t.Errorf("非结构体类型应返回空, got %d", len(got))
	}
}

func TestFieldsOfCached(t *testing.T) {
	a := FieldsOf(reflect.TypeFor[metaModel]())
	b := FieldsOf(reflect.TypeFor[metaModel]())
	if len(a) == 0 {
		t.Fatal("字段列表为空")
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("缓存结果不一致: %+v vs %+v", a[i], b[i])
		}
	}
}
