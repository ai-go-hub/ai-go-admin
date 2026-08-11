package util

import "testing"

// ==================== Str 相关测试 ====================

func TestSnakeToPascal(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"user_log", "UserLog"},
		{"auth_rule", "AuthRule"},
		{"auth/rule", "AuthRule"},
		{"123", "123"},
		{"a1_b2", "A1B2"},
		{"user", "User"},
		{"", ""},
		{"a__b", "AB"},
		{"__test__", "Test"},
	}
	for _, c := range cases {
		if got := SnakeToPascal(c.in); got != c.want {
			t.Errorf("SnakeToPascal(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestSnakeToLowerCamel(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"user_log", "userLog"},
		{"user_id", "userId"},
		{"auth_rule", "authRule"},
		{"auth/rule", "authRule"},
		{"created_at", "createdAt"},
		{"updated_at", "updatedAt"},
		{"start_time", "startTime"},
		{"a1_b2", "a1B2"},
		{"a_b_c", "aBC"},
		{"user", "user"},
		{"weigh", "weigh"},
		{"123", "123"},
		{"", ""},
		{"a__b", "aB"},
		{"__test__", "test"},
	}
	for _, c := range cases {
		if got := SnakeToLowerCamel(c.in); got != c.want {
			t.Errorf("SnakeToLowerCamel(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestPascalToSnake(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Admin", "admin"},
		{"AdminRule", "admin_rule"},
		{"AdminID", "admin_id"},
		{"UserLog", "user_log"},
		{"ID", "id"},
		{"APIKey", "api_key"},
		{"User2API", "user2_api"},
		{"UserAPI", "user_api"},
		{"中文", "中文"},
		{"123", "123"},
		{"", ""},
	}
	for _, c := range cases {
		if got := PascalToSnake(c.in); got != c.want {
			t.Errorf("PascalToSnake(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestTruncateStr(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want string
	}{
		{"英文短于 n 原样返回", "hello", 10, "hello"},
		{"英文等于 n 原样返回", "hello", 5, "hello"},
		{"英文长于 n 截断", "hello world", 5, "hello"},
		{"空串", "", 5, ""},
		{"n 为 0 截断为空串", "hello", 0, ""},
		{"中文按 rune 截断", "你好世界", 2, "你好"},
		{"中英混合按 rune 截断", "hi你好", 3, "hi你"},
		{"中文长度足够原样返回", "你好", 5, "你好"},
		{"emoji 按 rune 截断", "😀😃😄", 2, "😀😃"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateStr(tt.s, tt.n); got != tt.want {
				t.Errorf("TruncateStr(%q, %d) = %q, want %q", tt.s, tt.n, got, tt.want)
			}
		})
	}
}

// ==================== Ptr 相关测试 ====================

func TestFromPtr(t *testing.T) {
	t.Run("nil int 指针返回零值", func(t *testing.T) {
		var p *int
		if got := FromPtr(p); got != 0 {
			t.Errorf("FromPtr(nil) = %d, want 0", got)
		}
	})

	t.Run("非 nil int 指针返回原值", func(t *testing.T) {
		v := 42
		if got := FromPtr(&v); got != 42 {
			t.Errorf("FromPtr(&42) = %d, want 42", got)
		}
	})

	t.Run("nil string 指针返回空串", func(t *testing.T) {
		var p *string
		if got := FromPtr(p); got != "" {
			t.Errorf("FromPtr(nil) = %q, want \"\"", got)
		}
	})

	t.Run("非 nil string 指针返回原值", func(t *testing.T) {
		v := "hello"
		if got := FromPtr(&v); got != "hello" {
			t.Errorf("FromPtr(&\"hello\") = %q, want \"hello\"", got)
		}
	})

	t.Run("非 nil 但值为零值", func(t *testing.T) {
		v := 0
		if got := FromPtr(&v); got != 0 {
			t.Errorf("FromPtr(&0) = %d, want 0", got)
		}
	})

	t.Run("bool 指针", func(t *testing.T) {
		v := true
		if got := FromPtr(&v); got != true {
			t.Errorf("FromPtr(&true) = %v, want true", got)
		}
		var p *bool
		if got := FromPtr(p); got != false {
			t.Errorf("FromPtr(nil bool) = %v, want false", got)
		}
	})
}

func TestPtrIsZero(t *testing.T) {
	t.Run("nil 指针为 true", func(t *testing.T) {
		var p *int
		if !PtrIsZero(p) {
			t.Errorf("PtrIsZero(nil) = false, want true")
		}
	})

	t.Run("指向零值为 true", func(t *testing.T) {
		v := 0
		if !PtrIsZero(&v) {
			t.Errorf("PtrIsZero(&0) = false, want true")
		}
	})

	t.Run("指向非零值为 false", func(t *testing.T) {
		v := 1
		if PtrIsZero(&v) {
			t.Errorf("PtrIsZero(&1) = true, want false")
		}
	})

	t.Run("string 空串为 true", func(t *testing.T) {
		v := ""
		if !PtrIsZero(&v) {
			t.Errorf("PtrIsZero(&\"\") = false, want true")
		}
	})

	t.Run("string 非空为 false", func(t *testing.T) {
		v := "x"
		if PtrIsZero(&v) {
			t.Errorf("PtrIsZero(&\"x\") = true, want false")
		}
	})

	t.Run("bool false 视为零值", func(t *testing.T) {
		v := false
		if !PtrIsZero(&v) {
			t.Errorf("PtrIsZero(&false) = false, want true")
		}
	})
}

func TestPtrNotZero(t *testing.T) {
	t.Run("nil 指针为 false", func(t *testing.T) {
		var p *int
		if PtrNotZero(p) {
			t.Errorf("PtrNotZero(nil) = true, want false")
		}
	})

	t.Run("指向零值为 false", func(t *testing.T) {
		v := 0
		if PtrNotZero(&v) {
			t.Errorf("PtrNotZero(&0) = true, want false")
		}
	})

	t.Run("指向非零值为 true", func(t *testing.T) {
		v := 5
		if !PtrNotZero(&v) {
			t.Errorf("PtrNotZero(&5) = false, want true")
		}
	})

	t.Run("string 非空为 true", func(t *testing.T) {
		v := "hi"
		if !PtrNotZero(&v) {
			t.Errorf("PtrNotZero(&\"hi\") = false, want true")
		}
	})

	t.Run("bool true 为 true", func(t *testing.T) {
		v := true
		if !PtrNotZero(&v) {
			t.Errorf("PtrNotZero(&true) = false, want true")
		}
	})
}

func TestToPtr(t *testing.T) {
	t.Run("int 值转指针", func(t *testing.T) {
		p := ToPtr(123)
		if p == nil {
			t.Fatal("ToPtr(123) = nil, want non-nil")
		}
		if *p != 123 {
			t.Errorf("*ToPtr(123) = %d, want 123", *p)
		}
	})

	t.Run("string 值转指针", func(t *testing.T) {
		p := ToPtr("abc")
		if p == nil {
			t.Fatal("ToPtr(\"abc\") = nil, want non-nil")
		}
		if *p != "abc" {
			t.Errorf("*ToPtr(\"abc\") = %q, want \"abc\"", *p)
		}
	})

	t.Run("零值也返回非 nil 指针", func(t *testing.T) {
		p := ToPtr(0)
		if p == nil {
			t.Fatal("ToPtr(0) = nil, want non-nil")
		}
		if *p != 0 {
			t.Errorf("*ToPtr(0) = %d, want 0", *p)
		}
	})

	t.Run("struct 值转指针", func(t *testing.T) {
		type user struct{ Name string }
		p := ToPtr(user{Name: "yang"})
		if p == nil {
			t.Fatal("ToPtr(struct) = nil, want non-nil")
		}
		if p.Name != "yang" {
			t.Errorf("p.Name = %q, want \"yang\"", p.Name)
		}
	})

	t.Run("每次调用返回独立指针", func(t *testing.T) {
		p1 := ToPtr(1)
		p2 := ToPtr(1)
		if p1 == p2 {
			t.Errorf("ToPtr 两次调用返回了同一个指针")
		}
	})
}
