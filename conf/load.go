package conf

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/constraints"
	"gopkg.in/yaml.v3"
)

func init() {
	initFromYAML()
	mustInit()
}

func loadLocalFile(fname string) []byte {
	file, err := os.OpenFile(fname, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return b
}

var (
	confs map[string]any = map[string]any{}
)

var ErrConfigNotFound error = errors.New("configs not found")

func initFromYAML() {
	fname := "configs.yaml"
	if idx := slices.Index(os.Args, "-c"); idx != -1 {
		fname = os.Args[idx+1]
	}
	b := loadLocalFile(fname)
	if b == nil {
		return
	}
	err := yaml.Unmarshal(b, &confs)
	if err != nil {
		panic(fmt.Errorf("LoadFromJSON unmarshal error %v", err))
	}
}

// func initFromJSONC() {
// 	fname := "configs.jsonc"
// 	if idx := slices.Index(os.Args, "-c"); idx != -1 {
// 		fname = os.Args[idx+1]
// 	}
// 	b := loadLocalFile(fname)
// 	if b == nil {
// 		return
// 	}
// 	reg := regexp.MustCompile(`/\*{1,2}[\s\S]*?\*/`)
// 	b = reg.ReplaceAll(b, []byte("\n"))
// 	reg = regexp.MustCompile(`\s//[\s\S]*?\n`)
// 	b = reg.ReplaceAll(b, []byte("\n"))
// 	err := json.Unmarshal(b, &confs)
// 	if err != nil {
// 		panic(fmt.Errorf("LoadFromJSON unmarshal error %v", err))
// 	}
// }

func Any(name string) (v any, ok bool) {
	path := strings.Split(name, ".")
	var next any = confs
	for _, n := range path {
		tmp, ok := next.(map[string]any)
		if !ok {
			return v, false
		}
		next, ok = tmp[n]
		if !ok {
			return v, false
		}
	}
	// 类型错误触发panic中断
	return next, true
}

func Num[T constraints.Integer | constraints.Float](name string, default_ ...T) T {
	v, ok := Any(name)
	if ok {
		switch num := v.(type) {
		case int:
			return T(num)
		case float64:
			return T(num)
		default:
			panic("type error")
		}
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func Str(name string, default_ ...string) string {
	v, ok := Any(name)
	if ok {
		return v.(string)
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func Bool(name string, default_ ...bool) bool {
	v, ok := Any(name)
	if ok {
		return v.(bool)
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func List[T any](name string, default_ ...[]T) []T {
	a, ok := Any(name)
	if ok {
		v, ok := a.([]T)
		if ok {
			return v
		}
		var result []T
		mapstructure.Decode(a, &result)
		return result

	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func Map[T any](name string, default_ ...map[string]T) map[string]T {
	a, ok := Any(name)
	if ok {
		var m map[string]T
		err := mapstructure.Decode(a, &m)
		if err != nil {
			panic(err)
		}
		return m
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func Scan(name string, v any) error {
	a, ok := Any(name)
	if !ok {
		return ErrConfigNotFound
	}
	return mapstructure.Decode(a, v)
}
