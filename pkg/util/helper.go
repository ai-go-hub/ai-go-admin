package util

import "reflect"

// IsZero 判断 v 是否为零值或空值。
// 字符串额外把 "0" 视为零值，reflect 包对字符串零值检查方式是 v.Len() == 0
func IsZero(v any) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok {
		return s == "" || s == "0"
	}
	return reflect.ValueOf(v).IsZero()
}

// TruncateString 截断字符串，最多保留前 n 个字符（按 rune 截断，兼容多字节字符如中文）
func TruncateString(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}
