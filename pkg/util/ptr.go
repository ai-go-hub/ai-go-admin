package util

// FromPtr 安全解引用，nil 返回零值
func FromPtr[T comparable](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

// PtrIsZero 值为 nil 或值为对应类型零值
func PtrIsZero[T comparable](v *T) bool {
	return v == nil || FromPtr(v) == *new(T)
}

// PtrNotZero 值非 nil 并且值非对应类型零值
func PtrNotZero[T comparable](v *T) bool {
	return v != nil && FromPtr(v) != *new(T)
}
