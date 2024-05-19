package utils

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

// 格式化输出任意类型字符串
// 一般仅用在日志输出等需要有较强可读性的场景
func String(v any) string {
	type IString interface{ String() string }
	if v0, ok := v.(IString); ok {
		return fmt.Sprintf("[ %s [ %s ] ]", TypeName(v0), v0.String())
	}
	if b, err := json.Marshal(v); err == nil {
		return fmt.Sprintf("[ %s %s ]", TypeName(v), string(b))
	}
	return fmt.Sprintf("[ %s: %v ]", TypeName(v), v)
}

// 零拷贝字节数组转字符串
// 如果不确定是否存在并发则不使用
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 零拷贝字符串转字节数组
// 如果不确定是否存在并发则不使用
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
