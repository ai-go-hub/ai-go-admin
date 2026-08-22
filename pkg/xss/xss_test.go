package xss

import (
	"reflect"
	"strings"
	"testing"
)

// richModel 测试用模型: title=纯文本、content=富文本、note=富文本(*string)、remark=无tag、secret=json 隐藏
type richModel struct {
	ID      uint    `json:"id"`
	Title   string  `json:"title" xss:"text"`
	Content string  `json:"content" xss:"html"`
	Remark  *string `json:"remark"`
	Note    *string `json:"note" xss:"html"`
	Secret  string  `json:"-"`
}

func strPtr(s string) *string { return &s }

func TestHTMLPolicySanitize(t *testing.T) {
	got := HTMLPolicySanitize(`<p onclick="alert(1)">hello <b>world</b></p><script>alert(1)</script><a href="javascript:alert(1)">x</a>`)
	for _, bad := range []string{"<script", "onclick", "javascript:"} {
		if strings.Contains(got, bad) {
			t.Errorf("HTMLPolicySanitize 仍包含 %q: %s", bad, got)
		}
	}
	if !strings.Contains(got, "<b>world</b>") {
		t.Errorf("HTMLPolicySanitize 丢失合法标签: %s", got)
	}
	for _, keep := range []string{"hello", "x"} {
		if !strings.Contains(got, keep) {
			t.Errorf("HTMLPolicySanitize 丢失文本 %q: %s", keep, got)
		}
	}
}

func TestTextPolicySanitize(t *testing.T) {
	got := TextPolicySanitize(`<p>hello <b>world</b></p><script>alert(1)</script>`)
	if strings.ContainsAny(got, "<>") {
		t.Errorf("TextPolicySanitize 仍包含标签: %q", got)
	}
	if !strings.Contains(got, "hello world") {
		t.Errorf("TextPolicySanitize 丢失文本: %q", got)
	}
}

func TestPoliciesOf(t *testing.T) {
	m := policiesOf(reflect.TypeFor[richModel]())
	cases := map[string]Policy{
		"title":   PolicyText,
		"content": PolicyHTML,
		"note":    PolicyHTML,
	}
	for key, want := range cases {
		if m[key] != want {
			t.Errorf("policiesOf[%q] = %v, want %v", key, m[key], want)
		}
	}
	for _, key := range []string{"id", "remark", "secret"} {
		if _, ok := m[key]; ok {
			t.Errorf("policiesOf 不应包含 %q", key)
		}
	}
}

func TestSanitizeStruct(t *testing.T) {
	remark := `<b>ok</b>`
	m := richModel{
		Title:   `<p>hello</p>`,
		Content: `<p onclick="x">hi</p><script>alert(1)</script>`,
		Remark:  &remark,
		Note:    strPtr(`<a href="javascript:alert(1)">x</a>`),
		Secret:  `<script>alert(1)</script>`,
	}
	if err := SanitizeStruct(&m); err != nil {
		t.Fatal(err)
	}

	if m.Title != "hello" {
		t.Errorf("Title = %q, want hello", m.Title)
	}
	if strings.Contains(m.Content, "<script") || strings.Contains(m.Content, "onclick") {
		t.Errorf("Content 未清洗: %q", m.Content)
	}
	if m.Remark == nil || *m.Remark != "<b>ok</b>" {
		t.Errorf("Remark(无tag) 不应被修改: %v", m.Remark)
	}
	if m.Note == nil || strings.Contains(*m.Note, "javascript:") {
		t.Errorf("Note 未清洗: %v", m.Note)
	}
	if m.Secret != `<script>alert(1)</script>` {
		t.Errorf("Secret(json-) 不应被修改: %q", m.Secret)
	}
}

func TestSanitizeStructNested(t *testing.T) {
	type inner struct {
		Body string `json:"body" xss:"html"`
	}
	type outer struct {
		Name   string  `json:"name" xss:"text"`
		Single inner   `json:"single"`
		Items  []inner `json:"items"`
	}

	o := outer{
		Name:   `<b>n</b>`,
		Single: inner{Body: `<p onclick="a">s</p>`},
		Items:  []inner{{Body: `<script>alert(1)</script>`}, {Body: `<i>t</i>`}},
	}
	if err := SanitizeStruct(&o); err != nil {
		t.Fatal(err)
	}

	if o.Name != "n" {
		t.Errorf("Name = %q, want n", o.Name)
	}
	if strings.Contains(o.Single.Body, "onclick") || !strings.Contains(o.Single.Body, "s") {
		t.Errorf("Single.Body 未清洗: %q", o.Single.Body)
	}
	if len(o.Items) != 2 || strings.Contains(o.Items[0].Body, "<script") || o.Items[1].Body != "<i>t</i>" {
		t.Errorf("Items 未清洗: %+v", o.Items)
	}
}

func TestSanitizeStructNestedMapValue(t *testing.T) {
	type inner struct {
		Body string `json:"body" xss:"html"`
	}
	type model struct {
		Data map[string]inner `json:"data"`
	}
	v := model{Data: map[string]inner{"a": {Body: `<p onclick="x">hi</p>`}}}
	if err := SanitizeStruct(&v); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(v.Data["a"].Body, "onclick") || !strings.Contains(v.Data["a"].Body, "hi") {
		t.Errorf("map 值结构体字段未清洗: %q", v.Data["a"].Body)
	}
}

func TestSanitizeStructRejectNonStruct(t *testing.T) {
	if err := SanitizeStruct([]string{"<p>x</p>"}); err == nil {
		t.Error("传入 slice 应报错")
	}
	if err := SanitizeStruct(map[string]string{"a": "<p>x</p>"}); err == nil {
		t.Error("传入 map 应报错")
	}
	if err := SanitizeStruct(123); err == nil {
		t.Error("传入非结构体应报错")
	}
}

func TestSanitizeStructSliceStrings(t *testing.T) {
	type model struct {
		Items []string `json:"items" xss:"html"`
	}
	m := model{Items: []string{`<p onclick="x">a</p>`, `<script>alert(1)</script>`, `<b>c</b>`}}
	if err := SanitizeStruct(&m); err != nil {
		t.Fatal(err)
	}
	if len(m.Items) != 3 {
		t.Fatalf("len = %d, want 3", len(m.Items))
	}
	if strings.Contains(m.Items[0], "onclick") || strings.Contains(m.Items[1], "<script") {
		t.Errorf("slice 元素未清洗: %v", m.Items)
	}
	if m.Items[2] != "<b>c</b>" {
		t.Errorf("合法标签应保留: %q", m.Items[2])
	}
}

func TestSanitizeStructMapStrings(t *testing.T) {
	type model struct {
		Data map[string]string `json:"data" xss:"html"`
	}
	m := model{Data: map[string]string{
		"a": `<p onclick="x">a</p>`,
		"b": `<script>alert(1)</script>`,
	}}
	if err := SanitizeStruct(&m); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(m.Data["a"], "onclick") || strings.Contains(m.Data["b"], "<script") {
		t.Errorf("map value 未清洗: %v", m.Data)
	}
	if !strings.Contains(m.Data["a"], "a") {
		t.Errorf("map value 文本丢失: %q", m.Data["a"])
	}
}

// TestPoliciesOfSelfRefer 自引用类型不应死循环
func TestPoliciesOfSelfRefer(t *testing.T) {
	type node struct {
		Title string  `json:"title" xss:"text"`
		Kids  []*node `json:"kids"`
	}
	m := policiesOf(reflect.TypeFor[node]())
	if m["title"] != PolicyText {
		t.Errorf("title policy = %v, want PolicyText", m["title"])
	}
}

func TestSanitizeStructRejectTagOnSliceStruct(t *testing.T) {
	type inner struct {
		Body string `json:"body"`
	}
	type model struct {
		Blocks []inner `json:"blocks" xss:"html"`
	}
	if err := SanitizeStruct(&model{Blocks: []inner{{Body: "x"}}}); err == nil {
		t.Error("结构体切片字段声明清洗 tag 应报错")
	}
}

func TestSanitizeStructRejectTagOnMapStruct(t *testing.T) {
	type inner struct {
		Body string `json:"body"`
	}
	type model struct {
		Data map[string]inner `json:"data" xss:"html"`
	}
	if err := SanitizeStruct(&model{Data: map[string]inner{"a": {Body: "x"}}}); err == nil {
		t.Error("结构体 map 字段声明清洗 tag 应报错")
	}
}
