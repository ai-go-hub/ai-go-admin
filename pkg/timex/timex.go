package timex

import "time"

// Format 将 time.Time 按指定 layout 格式化为字符串，零值返回空字符串
func Format(t time.Time, layout string) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(layout)
}
