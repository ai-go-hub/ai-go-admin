package util

// TruncateStr 截断字符串，最多保留前 n 个字符（按 rune 截断，兼容多字节字符如中文）
func TruncateStr(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}
