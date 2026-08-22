// Package xss 提供基于 HTML 白名单策略的反 XSS 清洗能力
//
// 通过结构体字段上的 xss tag 声明清洗策略:
//
//	xss:"html" -> HTMLPolicySanitize: 富文本白名单清洗，保留常用标签，去除 script/事件属性/javascript: 等
//	xss:"text" -> TextPolicySanitize: 剥离全部 HTML 标签，得到纯文本
//	无 xss tag  -> 不过滤
package xss

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/ai-go-hub/ai-go-admin/pkg/structx"
	"github.com/microcosm-cc/bluemonday"
)

// Policy 清洗策略
type Policy uint8

const (
	// PolicyNone 不过滤（默认）
	PolicyNone Policy = iota
	// PolicyHTML 富文本白名单清洗
	PolicyHTML
	// PolicyText 剥离全部 HTML 标签
	PolicyText
)

var (
	// htmlPolicy 富文本白名单策略
	htmlPolicy = bluemonday.UGCPolicy()
	// textPolicy 纯文本严格策略
	textPolicy = bluemonday.StrictPolicy()
)

// HTMLPolicySanitize 富文本过滤: 白名单保留常用标签，去除 script/事件属性/javascript: 协议等
func HTMLPolicySanitize(s string) string {
	return htmlPolicy.Sanitize(s)
}

// TextPolicySanitize 剥离全部 HTML 标签，返回纯文本
func TextPolicySanitize(s string) string {
	return textPolicy.Sanitize(s)
}

// Apply 按策略清洗单个字符串
func Apply(p Policy, s string) string {
	switch p {
	case PolicyHTML:
		return HTMLPolicySanitize(s)
	case PolicyText:
		return TextPolicySanitize(s)
	default:
		return s
	}
}

// SanitizeStruct 按结构体字段的 xss tag 清洗 string / *string 值
// 支持嵌套结构体，根据其字段的 xss tag 正常清洗
// 字符串容器（[]string、map[string]string 等）遍历清洗
// 传入 slice/map 等非结构体返回错误（防止以为洗了，实际没洗的安全隐患）
func SanitizeStruct(v any) error {
	if v == nil {
		return nil
	}
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("xss: 只支持清洗结构体，收到类型 %s", t)
	}
	// 快速跳过: 没有任何 xss tag
	if len(policiesOf(reflect.TypeOf(v))) == 0 {
		return nil
	}
	return sanitizeValue(reflect.ValueOf(v))
}

// ==================== 以下为内部实现 ====================

// policiesCache 类型 -> jsonKey -> Policy 缓存
var policiesCache sync.Map

// policiesOf 返回类型中带 xss tag 字段的 jsonKey -> Policy 映射
func policiesOf(t reflect.Type) map[string]Policy {
	if cached, ok := policiesCache.Load(t); ok {
		return cached.(map[string]Policy)
	}

	m := make(map[string]Policy)
	collectPolicies(t, m, make(map[reflect.Type]bool))
	actual, _ := policiesCache.LoadOrStore(t, m)
	return actual.(map[string]Policy)
}

// policyOf 解析字段 xss tag，返回清洗策略
func policyOf(tag string) Policy {
	name := strings.TrimSpace(strings.SplitN(tag, ",", 2)[0])
	switch name {
	case "html":
		return PolicyHTML
	case "text":
		return PolicyText
	default:
		return PolicyNone
	}
}

// collectPolicies 递归收集类型内所有带 xss tag 字段的策略（visited 防止自引用类型死循环）
func collectPolicies(t reflect.Type, m map[string]Policy, visited map[reflect.Type]bool) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if visited[t] {
		return
	}

	switch t.Kind() {
	case reflect.Struct:
		visited[t] = true
		for _, f := range structx.FieldsOf(t) {
			if p := policyOf(f.Tag.Get("xss")); p != PolicyNone {
				if f.Key != "" {
					m[f.Key] = p
				}
			} else {
				collectPolicies(f.Type, m, visited)
			}
		}
	case reflect.Slice:
		collectPolicies(t.Elem(), m, visited)
	case reflect.Map:
		collectPolicies(t.Elem(), m, visited)
	}
}

// sanitizeValue 递归清洗
func sanitizeValue(rv reflect.Value) error {
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		for _, f := range structx.FieldsOf(rv.Type()) {
			fv := rv.Field(f.Index)
			p := policyOf(f.Tag.Get("xss"))
			if p == PolicyNone {
				if err := sanitizeValue(fv); err != nil {
					return err
				}
				continue
			}
			if err := sanitizeField(fv, p, f.Key); err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			if err := sanitizeValue(rv.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, k := range rv.MapKeys() {
			mv := rv.MapIndex(k)
			if !mv.CanInterface() {
				continue
			}
			// map 值不可寻址，复制一份递归清洗后写回
			nv := reflect.New(mv.Type()).Elem()
			nv.Set(mv)
			if err := sanitizeValue(nv); err != nil {
				return err
			}
			rv.SetMapIndex(k, nv)
		}
	}
	return nil
}

// sanitizeField 按策略清洗带 tag 的字段值
// []string/[]*string 与 map[string]string/map[string]*string 遍历元素清洗
func sanitizeField(fv reflect.Value, p Policy, key string) error {
	switch fv.Kind() {
	case reflect.String:
		if fv.CanSet() {
			fv.SetString(Apply(p, fv.String()))
		}
	case reflect.Pointer:
		if fv.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("xss: 字段 %q 为 %s 类型，不支持声明清洗 tag", key, fv.Type().Elem().Kind())
		}
		if !fv.IsNil() && fv.Elem().CanSet() {
			fv.Elem().SetString(Apply(p, fv.Elem().String()))
		}
	case reflect.Slice:
		elem := fv.Type().Elem()
		switch {
		case elem.Kind() == reflect.String:
			for i := 0; i < fv.Len(); i++ {
				if e := fv.Index(i); e.CanSet() {
					e.SetString(Apply(p, e.String()))
				}
			}
		case elem.Kind() == reflect.Pointer && elem.Elem().Kind() == reflect.String:
			for i := 0; i < fv.Len(); i++ {
				if e := fv.Index(i); !e.IsNil() && e.Elem().CanSet() {
					e.Elem().SetString(Apply(p, e.Elem().String()))
				}
			}
		default:
			return fmt.Errorf("xss: 字段 %q 为 %s 类型，不支持声明清洗 tag", key, elem)
		}
	case reflect.Map:
		vt := fv.Type().Elem()
		switch {
		case vt.Kind() == reflect.String:
			for _, k := range fv.MapKeys() {
				nv := reflect.New(vt).Elem()
				nv.Set(fv.MapIndex(k))
				nv.SetString(Apply(p, nv.String()))
				fv.SetMapIndex(k, nv)
			}
		case vt.Kind() == reflect.Pointer && vt.Elem().Kind() == reflect.String:
			for _, k := range fv.MapKeys() {
				nv := reflect.New(vt).Elem()
				nv.Set(fv.MapIndex(k))
				if !nv.IsNil() && nv.Elem().CanSet() {
					nv.Elem().SetString(Apply(p, nv.Elem().String()))
				}
				fv.SetMapIndex(k, nv)
			}
		default:
			return fmt.Errorf("xss: 字段 %q 为 %s 类型，不支持声明清洗 tag", key, vt)
		}
	default:
		return fmt.Errorf("xss: 字段 %q 为 %s 类型，不支持声明清洗 tag", key, fv.Kind())
	}
	return nil
}
