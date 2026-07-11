package copierx

import (
	"time"

	"github.com/jinzhu/copier"
)

// Time 返回一个 copier.TypeConverter，将 time.Time 按指定 layout 格式化为字符串
// 零值时间（IsZero() == true）返回空字符串
func Time(layout string) copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: time.Time{},
		DstType: copier.String,
		Fn: func(src any) (any, error) {
			t := src.(time.Time)
			if t.IsZero() {
				return "", nil
			}
			return t.Format(layout), nil
		},
	}
}
