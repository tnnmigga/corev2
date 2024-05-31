package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

func Init(b []byte) {
	LoadFromJSON(b)
	serverID = Int("server.id")
	serverType = String("server.type")
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
