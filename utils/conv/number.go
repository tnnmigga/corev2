package conv

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

// 将任意类型值转换为指定类型整型
func Integer[T constraints.Integer](value any) T {
	switch v := value.(type) {
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		return T(n)
	case int:
		return T(v)
	case int16:
		return T(v)
	case int32:
		return T(v)
	case int64:
		return T(v)
	case uint:
		return T(v)
	case uint16:
		return T(v)
	case uint32:
		return T(v)
	case uint64:
		return T(v)
	case float32:
		return T(v)
	case float64:
		return T(v)
	case uintptr:
		return T(v)
	case byte:
		return T(v)
	default:
		panic(fmt.Errorf("Integer transfor error %#v", v))
	}
}

// 将任意类型值转换为指定类型浮点型
func Float[T constraints.Float](value any) T {
	switch v := value.(type) {
	case string:
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(err)
		}
		return T(n)
	case int:
		return T(v)
	case int16:
		return T(v)
	case int32:
		return T(v)
	case int64:
		return T(v)
	case uint:
		return T(v)
	case uint16:
		return T(v)
	case uint32:
		return T(v)
	case uint64:
		return T(v)
	case float32:
		return T(v)
	case float64:
		return T(v)
	default:
		panic(fmt.Errorf("Float type error %#v", value))
	}
}
