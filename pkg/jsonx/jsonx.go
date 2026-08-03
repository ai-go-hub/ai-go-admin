// Package jsonx 提供泛型 JSON 解码
package jsonx

import "encoding/json"

// UnmarshalSafe 泛型 JSON 解码，解析失败静默返回零值
func UnmarshalSafe[T any](data []byte) T {
	var v T
	json.Unmarshal(data, &v)
	return v
}
