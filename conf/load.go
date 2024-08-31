package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func init() {
	fname := "configs.jsonc"
	b := loadLocalFile(fname)
	if b == nil {
		return
	}
	initFromJSON(b)
	mustInit()
}

func loadLocalFile(fname string) []byte {
	file, err := os.OpenFile(fname, os.O_RDONLY, 0)
	if err != nil {
		log.Println(err)
		return nil
	}
	b, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return nil
	}
	return b
}

var (
	confs map[string]any = map[string]any{}
)

var ErrConfigNotFound error = errors.New("configs not found")

func initFromJSON(b []byte) {
	b = uncomment(b)
	err := json.Unmarshal(b, &confs)
	if err != nil {
		panic(fmt.Errorf("LoadFromJSON unmarshal error %v", err))
	}
}

func uncomment(b []byte) []byte {
	reg := regexp.MustCompile(`/\*{1,2}[\s\S]*?\*/`)
	b = reg.ReplaceAll(b, []byte("\n"))
	reg = regexp.MustCompile(`\s//[\s\S]*?\n`)
	return reg.ReplaceAll(b, []byte("\n"))
}

func Any[T any](name string) (v T, ok bool) {
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
	return next.(T), true
}

func Int(name string, default_ ...int) int {
	v, ok := Any[float64](name)
	if ok {
		return int(v)
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func Int64(name string, default_ ...int64) int64 {
	v, ok := Any[float64](name)
	if ok {
		return int64(v)
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func Int32(name string, default_ ...int32) int32 {
	v, ok := Any[float64](name)
	if ok {
		return int32(v)
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func Uint64(name string, default_ ...uint64) uint64 {
	v, ok := Any[float64](name)
	if ok {
		return uint64(v)
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func Uint32(name string, default_ ...uint32) uint32 {
	v, ok := Any[float64](name)
	if ok {
		return uint32(v)
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func String(name string, default_ ...string) string {
	v, ok := Any[string](name)
	if ok {
		return string(v)
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func Float64(name string, default_ ...float64) float64 {
	v, ok := Any[float64](name)
	if ok {
		return v
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func Bool(name string, default_ ...bool) bool {
	v, ok := Any[bool](name)
	if ok {
		return v
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func List[T any](name string, default_ ...[]T) []T {
	a, ok := Any[[]any](name)
	if ok {
		ar := make([]T, len(a))
		for i, v := range a {
			ar[i] = v.(T)
		}
		return ar
	}
	if len(default_) > 0 {
		return default_[0]
	}
	panic(ErrConfigNotFound)
}

func Map[T any](name string, default_ ...map[string]T) map[string]T {
	a, ok := Any[map[string]any](name)
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
	a, ok := Any[any](name)
	if !ok {
		return ErrConfigNotFound
	}
	return mapstructure.Decode(a, v)
}
