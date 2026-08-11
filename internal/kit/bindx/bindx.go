package bindx

import (
	"encoding/json"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/copier"
)

// Tri 三绑定结果容器
type Tri[T any] struct {
	Model T
	Map   map[string]any
	DTO   any
}

// ShouldBindTri 三绑定，将单份原始请求体同时绑定至多个对象
//
//   - 绑定与校验目标为 DTO，未设置 DTO 时退化为绑定模型本身（适合无需 DTO 的场景）
//   - 若设置了 DTO，模型值对 DTO 使用 copier.Copy 转换而来，按字段名拷贝值
//   - Map 用于所见即所得的更新（排除结构体更新时的默认赋零值和自动跳过零值机制）
//   - Map 赋值以 DTO 的 json tag 为准，值已类型化（如 jsonb 的 *datatypes.JSON）
//
// 仅返回 err
func ShouldBindTri[T any](body []byte, tri *Tri[T]) error {
	target := tri.DTO
	if target == nil {
		target = &tri.Model
	}

	if err := binding.JSON.BindBody(body, target); err != nil {
		return err
	}

	// map 只需区分"传了 / 未传"，用 RawMessage 只解析 key、不解码值，成本最低
	var present map[string]json.RawMessage
	if err := json.Unmarshal(body, &present); err != nil {
		return err
	}

	tri.Map = StructToMap(reflect.ValueOf(target).Elem(), present)

	// DTO -> 模型
	if target != &tri.Model {
		if err := copier.Copy(&tri.Model, target); err != nil {
			return err
		}
	}
	return nil
}

// bindField 结构体字段元数据
type bindField struct {
	idx int
	key string
}

// bindFieldsCache 结构体字段缓存
// map[reflect.Type][]bindField
var bindFieldsCache sync.Map

// StructToMap 由结构体类型化的字段值构建 map
func StructToMap(rv reflect.Value, present map[string]json.RawMessage) map[string]any {
	rt := rv.Type()
	fields, ok := bindFieldsCache.Load(rt)
	if !ok {
		var list []bindField
		for i := 0; i < rt.NumField(); i++ {
			sf := rt.Field(i)
			if sf.PkgPath != "" {
				continue
			}
			key := sf.Tag.Get("json")
			if key == "" || key == "-" {
				continue
			}
			if idx := strings.Index(key, ","); idx >= 0 {
				key = key[:idx]
			}
			list = append(list, bindField{idx: i, key: key})
		}
		actual, _ := bindFieldsCache.LoadOrStore(rt, list)
		fields = actual
	}

	entity := make(map[string]any, len(present))
	for _, f := range fields.([]bindField) {
		if _, ok := present[f.key]; ok {
			entity[f.key] = rv.Field(f.idx).Interface()
		}
	}
	return entity
}
