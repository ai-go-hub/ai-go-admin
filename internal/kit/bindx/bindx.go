package bindx

import (
	"encoding/json"
	"reflect"

	"github.com/ai-go-hub/ai-go-admin/pkg/structx"
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

// StructToMap 使用结构体类型化的字段值构建 map
func StructToMap(rv reflect.Value, present map[string]json.RawMessage) map[string]any {
	entity := make(map[string]any, len(present))
	for _, f := range structx.FieldsOf(rv.Type()) {
		if f.Key == "" {
			continue
		}
		if _, ok := present[f.Key]; ok {
			entity[f.Key] = rv.Field(f.Index).Interface()
		}
	}
	return entity
}
