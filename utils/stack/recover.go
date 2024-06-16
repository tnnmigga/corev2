package stack

import (
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/tnnmigga/corev2/utils/idgen"
	"github.com/tnnmigga/corev2/zlog"
)

func RecoverPanic() {
	if r := recover(); r != nil {
		zlog.Errorf("%v: %s", r, debug.Stack())
	}
}

func ExecAndRecover(fn func()) {
	defer RecoverPanic()
	fn()
}

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

// 获取包名
func PkgName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	name := runtime.FuncForPC(pc).Name()
	return strings.Split(name, ".")[0]
}

// 获取结构体名称
func TypeName(v any) string {
	mType := reflect.TypeOf(v)
	for mType.Kind() == reflect.Ptr {
		mType = mType.Elem()
	}
	return mType.Name()
}

func TypeID(v any) uint32 {
	name := TypeName(v)
	return idgen.HashToID(name)
}

// 获取函数名称
func FuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
