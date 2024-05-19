package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/tnnmigga/corev2/proc"
)

func init() {
	fname := proc.Argv.Str("-c", "configs.jsonc")
	b := loadLocalFile(fname)
	LoadFromJSON(b)
	serverID = Int("server.id")
	serverType = String("server.type")
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

var errConfigNotFound error = errors.New("configs not found")

func LoadFromJSON(b []byte) {
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
