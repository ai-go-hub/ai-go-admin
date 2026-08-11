package bindx

import (
	"encoding/json"
	"reflect"
	"testing"

	"gorm.io/datatypes"
)

// ==================== Bind 相关测试 ====================

type bindTestModel struct {
	ID       uint            `json:"id"`
	Name     *string         `json:"name"`
	Bio      string          `json:"bio"`
	Status   string          `json:"status" binding:"required,oneof=enable disable"`
	AdminIDs *datatypes.JSON `json:"admin_ids"`
	Pwd      string          `json:"-"`
}

func TestShouldBindTri(t *testing.T) {
	body := []byte(`{"id":3,"name":null,"status":"enable","admin_ids":[2,3,4]}`)
	var tri Tri[bindTestModel]
	if err := ShouldBindTri(body, &tri); err != nil {
		t.Fatalf("ShouldBindTri err: %v", err)
	}

	if tri.Model.ID != 3 || tri.Model.Status != "enable" {
		t.Errorf("struct bound wrong: %+v", tri.Model)
	}
	if tri.Model.Name != nil {
		t.Errorf("name should stay nil, got %v", *tri.Model.Name)
	}

	// 未传字段不应出现在 map（跳过更新）
	if _, ok := tri.Map["bio"]; ok {
		t.Errorf("bio not sent but present in map: %v", tri.Map)
	}
	// json:"-" 字段不应出现
	if _, ok := tri.Map["pwd"]; ok {
		t.Errorf("json:\"-\" field present in map: %v", tri.Map)
	}
	// null 保留为 key，值类型化（*string 的 nil）
	ptr, ok := tri.Map["name"].(*string)
	if !ok || ptr != nil {
		t.Errorf("name should be present as nil *string, got %T %v", tri.Map["name"], tri.Map["name"])
	}
	// jsonb 类型化值
	v, ok := tri.Map["admin_ids"]
	if !ok {
		t.Fatalf("admin_ids missing: %v", tri.Map)
	}
	got, ok := v.(*datatypes.JSON)
	if !ok || string(*got) != "[2,3,4]" {
		t.Errorf("admin_ids not *datatypes.JSON([2,3,4]), got %T %v", v, v)
	}
}

func TestShouldBindTriValidation(t *testing.T) {
	var tri Tri[bindTestModel]
	if err := ShouldBindTri([]byte(`{"status":"bogus"}`), &tri); err == nil {
		t.Fatal("expected binding validation error")
	}
}

func TestShouldBindTriDropsUnknownKeys(t *testing.T) {
	var tri Tri[bindTestModel]
	if err := ShouldBindTri([]byte(`{"id":1,"status":"enable","hacker_key":123}`), &tri); err != nil {
		t.Fatalf("err: %v", err)
	}
	if _, ok := tri.Map["hacker_key"]; ok {
		t.Errorf("unknown key should be dropped: %v", tri.Map)
	}
}

func TestStructToMap(t *testing.T) {
	body := []byte(`{"id":3,"name":null,"status":"enable","admin_ids":[2,3,4]}`)
	var model bindTestModel
	if err := json.Unmarshal(body, &model); err != nil {
		t.Fatal(err)
	}
	var present map[string]json.RawMessage
	if err := json.Unmarshal(body, &present); err != nil {
		t.Fatal(err)
	}

	m := StructToMap(reflect.ValueOf(&model).Elem(), present)

	// 未传 / json:"-" / 未知 key 都不应出现
	if _, ok := m["bio"]; ok {
		t.Errorf("bio not sent but present in map: %v", m)
	}
	if _, ok := m["pwd"]; ok {
		t.Errorf("json:\"-\" field present in map: %v", m)
	}
	if _, ok := m["hacker_extra"]; ok {
		t.Errorf("unknown key present in map: %v", m)
	}
	// null 保留为 nil *string
	if ptr, ok := m["name"].(*string); !ok || ptr != nil {
		t.Errorf("name should be nil *string, got %T %v", m["name"], m["name"])
	}
	// jsonb 类型化值
	v, ok := m["admin_ids"]
	if !ok {
		t.Fatalf("admin_ids missing: %v", m)
	}
	got, ok := v.(*datatypes.JSON)
	if !ok || string(*got) != "[2,3,4]" {
		t.Errorf("admin_ids not *datatypes.JSON([2,3,4]), got %T %v", v, v)
	}
}

type bindUpdateDTO struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password"`
}

type bindUpdateModel struct {
	Name     string `json:"name"`
	Password string `json:"-"`
	Bio      string `json:"bio"`
}

// TestShouldBindTriWithDTO D != T 场景: DTO 负责绑定校验，模型 Model 由 copier 填充，Map 以 DTO 的 json tag 为准
func TestShouldBindTriWithDTO(t *testing.T) {
	body := []byte(`{"name":"x","password":"secret","bio":null}`)
	var tri Tri[bindUpdateModel]
	tri.DTO = &bindUpdateDTO{}
	if err := ShouldBindTri(body, &tri); err != nil {
		t.Fatalf("ShouldBindTri err: %v", err)
	}

	// 模型 Model 由 DTO 经 copier 填充（含 json:"-" 的密码字段）
	if tri.Model.Name != "x" || tri.Model.Password != "secret" {
		t.Errorf("model struct wrong: %+v", tri.Model)
	}
	// Map 以 DTO 的 json tag 为准（密码也在内）
	if v, ok := tri.Map["password"]; !ok || v != "secret" {
		t.Errorf("password not keyed by dto tag in map: %v", tri.Map)
	}
	// 未传字段不出现在 map（跳过更新）
	if _, ok := tri.Map["bio"]; ok {
		t.Errorf("bio not sent but present in map: %v", tri.Map)
	}

	// DTO 必填校验
	var tri2 Tri[bindUpdateModel]
	tri2.DTO = &bindUpdateDTO{}
	if err := ShouldBindTri([]byte(`{"password":"x"}`), &tri2); err == nil {
		t.Fatal("expected validation error (name required)")
	}
}
