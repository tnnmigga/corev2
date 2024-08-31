package utils

import (
	"reflect"
	"runtime"
)

// 获取调用者
func Caller(skip ...int) string {
	n := 1
	if len(skip) > 0 {
		n = skip[0]
	}
	pc, _, _, ok := runtime.Caller(n)
	if !ok {
		return "runtime.Caller() failed"
	}
	return runtime.FuncForPC(pc).Name()
}

// 获取结构体名称
func TypeName(v any) string {
	mType := reflect.TypeOf(v)
	for mType.Kind() == reflect.Ptr {
		mType = mType.Elem()
	}
	return mType.Name()
}

// 获取函数名称
func FuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
