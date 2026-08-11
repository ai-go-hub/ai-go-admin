package util

import (
	"strconv"
	"strings"
)

// TruncateStr 截断字符串，最多保留前 n 个字符（按 rune 截断，兼容多字节字符如中文）
func TruncateStr(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}

// SnakeToPascal 下划线转大驼峰，如 user_id -> UserId
func SnakeToPascal(s string) string {
	var b strings.Builder
	for _, seg := range strings.FieldsFunc(s, func(r rune) bool { return r == '_' || r == '/' }) {
		b.WriteString(strings.ToUpper(seg[:1]))
		b.WriteString(seg[1:])
	}
	return b.String()
}

// PascalToSnake 大驼峰转下划线小写，如 UserId -> user_id
func PascalToSnake(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				prev := s[i-1]
				nextLower := i+1 < len(s) && s[i+1] >= 'a' && s[i+1] <= 'z'
				if (prev >= 'a' && prev <= 'z') || (prev >= '0' && prev <= '9') || nextLower {
					b.WriteByte('_')
				}
			}
			b.WriteByte(c + 32)
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}

// SnakeToLowerCamel 下划线转小驼峰，如 user_id -> userId
func SnakeToLowerCamel(s string) string {
	var b strings.Builder
	for i, seg := range strings.FieldsFunc(s, func(r rune) bool { return r == '_' || r == '/' }) {
		if seg == "" {
			continue
		}
		if i > 0 {
			b.WriteString(strings.ToUpper(seg[:1]))
			b.WriteString(seg[1:])
		} else {
			b.WriteString(seg)
		}
	}
	return b.String()
}

// UintsToStrs uint 切片转字符串切片，如 [1,2] -> ["1","2"]
func UintsToStrs(ids []uint) []string {
	result := make([]string, 0, len(ids))
	for _, id := range ids {
		result = append(result, strconv.FormatUint(uint64(id), 10))
	}
	return result
}
