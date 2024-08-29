package env

import (
	"os"
	"strconv"
)

var Argv = argvParser{}

type argvParser struct {
}

// 查找命令行参数中的指定参数并解析成整型
func (a argvParser) Int(name string, default_ int) int {
	value := a.Str(name, "")
	if value == "" {
		return default_
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return n
}

// 查找命令行参数中的指定参数
func (a argvParser) Str(name string, default_ string) string {
	for i, v := range os.Args[1:] {
		if v == name {
			return os.Args[i+2]
		}
	}
	return default_
}

// 查找是否存在指定名称的命令行参数
func (a argvParser) Find(name string) bool {
	for _, v := range os.Args[1:] {
		if v == name {
			return true
		}
	}
	return false
}
