package tree

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestBuild(t *testing.T) {
	data := []map[string]any{
		{"id": 1, "pid": 0, "title": "标题1"},
		{"id": 2, "pid": 1, "title": "标题1-1"},
		{"id": 3, "pid": 1, "title": "标题1-2"},
		{"id": 4, "pid": 2, "title": "标题1-1-1"},
	}

	got := Build(data, "id", "pid", "children")

	want := []map[string]any{
		{
			"id":    1,
			"pid":   0,
			"title": "标题1",
			"children": []map[string]any{
				{
					"id":    2,
					"pid":   1,
					"title": "标题1-1",
					"children": []map[string]any{
						{"id": 4, "pid": 2, "title": "标题1-1-1"},
					},
				},
				{"id": 3, "pid": 1, "title": "标题1-2"},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Build() = %v, want %v", got, want)
	}
}

func TestBuildWithCustomChildrenField(t *testing.T) {
	data := []map[string]any{
		{"id": 1, "pid": 0, "title": "标题1"},
		{"id": 2, "pid": 1, "title": "标题1-1"},
	}

	got := Build(data, "id", "pid", "subs")

	want := []map[string]any{
		{
			"id":    1,
			"pid":   0,
			"title": "标题1",
			"subs": []map[string]any{
				{"id": 2, "pid": 1, "title": "标题1-1"},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Build() = %v, want %v", got, want)
	}
}

func TestRender(t *testing.T) {
	data := []map[string]any{
		{"id": 2, "pid": 0, "title": "权限管理", "name": "auth"},
		{"id": 3, "pid": 2, "title": "角色组管理", "name": "auth/group"},
		{"id": 8, "pid": 2, "title": "管理员管理", "name": "auth/admin"},
	}

	got := Render(data, "id", "pid", "title")

	want := []map[string]any{
		{"id": 2, "pid": 0, "title": "权限管理", "name": "auth"},
		{"id": 3, "pid": 2, "title": "    ├角色组管理", "name": "auth/group"},
		{"id": 8, "pid": 2, "title": "    └管理员管理", "name": "auth/admin"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render() = %v, want %v", got, want)
	}
}

func TestRenderUnordered(t *testing.T) {
	data := []map[string]any{
		{"id": 1, "pid": 0, "title": "A"},
		{"id": 2, "pid": 1, "title": "B"},
		{"id": 5, "pid": 0, "title": "E"},
		{"id": 6, "pid": 5, "title": "F"},
		{"id": 7, "pid": 5, "title": "G"},
		{"id": 8, "pid": 7, "title": "H"},
		{"id": 9, "pid": 7, "title": "I"},
		{"id": 3, "pid": 1, "title": "C"},
		{"id": 4, "pid": 1, "title": "D"},
	}

	got := Render(data, "id", "pid", "title")

	want := []map[string]any{
		{"id": 1, "pid": 0, "title": "A"},
		{"id": 2, "pid": 1, "title": "    ├B"},
		{"id": 3, "pid": 1, "title": "    ├C"},
		{"id": 4, "pid": 1, "title": "    └D"},
		{"id": 5, "pid": 0, "title": "E"},
		{"id": 6, "pid": 5, "title": "    ├F"},
		{"id": 7, "pid": 5, "title": "    └G"},
		{"id": 8, "pid": 7, "title": "         ├H"},
		{"id": 9, "pid": 7, "title": "         └I"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render() = %v, want %v", got, want)
	}
}

func TestRenderNested(t *testing.T) {
	data := []map[string]any{
		{"id": 1, "pid": 0, "title": "root"},
		{"id": 2, "pid": 1, "title": "a"},
		{"id": 3, "pid": 2, "title": "a1"},
		{"id": 4, "pid": 2, "title": "a2"},
		{"id": 5, "pid": 1, "title": "b"},
	}

	got := Render(data, "id", "pid", "title")

	want := []map[string]any{
		{"id": 1, "pid": 0, "title": "root"},
		{"id": 2, "pid": 1, "title": "    ├a"},
		{"id": 3, "pid": 2, "title": "    │    ├a1"},
		{"id": 4, "pid": 2, "title": "    │    └a2"},
		{"id": 5, "pid": 1, "title": "    └b"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render() = %v, want %v", got, want)
	}
}

func TestBuildEmpty(t *testing.T) {
	got := Build(nil, "id", "pid", "children")
	if got == nil {
		t.Error("Build(nil) should return empty slice, not nil")
	}
	if len(got) != 0 {
		t.Errorf("Build(nil) length = %d, want 0", len(got))
	}
}

func TestRenderEmpty(t *testing.T) {
	got := Render(nil, "id", "pid", "title")
	if got == nil {
		t.Error("Render(nil) should return empty slice, not nil")
	}
	if len(got) != 0 {
		t.Errorf("Render(nil) length = %d, want 0", len(got))
	}
}

// ==================== bench ====================

type benchRule struct {
	ID    uint   `json:"id"`
	Pid   *uint  `json:"pid"`
	Title string `json:"title"`
	Name  string `json:"name"`
	Path  string `json:"path"`
	Icon  string `json:"icon"`
	Type  string `json:"type"`
}

func makeRules(n int) []benchRule {
	rules := make([]benchRule, n)
	for i := range n {
		var pid *uint
		if i%5 != 0 {
			p := uint(i/5*5 + 1)
			pid = &p
		}
		rules[i] = benchRule{
			ID:    uint(i + 1),
			Pid:   pid,
			Title: "规则标题",
			Name:  "auth/rule",
			Path:  "/admin/rule",
			Icon:  "el-icon-menu",
			Type:  "menu",
		}
	}
	return rules
}

func BenchmarkManualConvert100(b *testing.B) {
	rules := makeRules(100)

	for b.Loop() {
		data := make([]map[string]any, len(rules))
		for j, r := range rules {
			pid := uint(0)
			if r.Pid != nil {
				pid = *r.Pid
			}
			data[j] = map[string]any{
				"id":    r.ID,
				"pid":   pid,
				"title": r.Title,
				"name":  r.Name,
				"path":  r.Path,
				"icon":  r.Icon,
				"type":  r.Type,
			}
		}
		_ = Build(data, "id", "pid", "children")
	}
}

func BenchmarkJSONConvert100(b *testing.B) {
	rules := makeRules(100)

	for b.Loop() {
		buf, _ := json.Marshal(rules)
		var data []map[string]any
		_ = json.Unmarshal(buf, &data)
		_ = Build(data, "id", "pid", "children")
	}
}

func BenchmarkManualConvert1000(b *testing.B) {
	rules := makeRules(1000)

	for b.Loop() {
		data := make([]map[string]any, len(rules))
		for j, r := range rules {
			pid := uint(0)
			if r.Pid != nil {
				pid = *r.Pid
			}
			data[j] = map[string]any{
				"id":    r.ID,
				"pid":   pid,
				"title": r.Title,
				"name":  r.Name,
				"path":  r.Path,
				"icon":  r.Icon,
				"type":  r.Type,
			}
		}
		_ = Build(data, "id", "pid", "children")
	}
}

func BenchmarkJSONConvert1000(b *testing.B) {
	rules := makeRules(1000)

	for b.Loop() {
		buf, _ := json.Marshal(rules)
		var data []map[string]any
		_ = json.Unmarshal(buf, &data)
		_ = Build(data, "id", "pid", "children")
	}
}
