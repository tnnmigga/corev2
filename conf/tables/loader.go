package tables

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type fromLocalJSON struct {
	path string
}

func (h fromLocalJSON) Load() (map[string]any, error) {
	result := map[string]any{}
	err := filepath.Walk(h.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var data any
		if err := json.Unmarshal(b, &data); err != nil {
			return err
		}
		result[info.Name()] = data
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func LoadFromLocalJSON(path string) ILoader {
	return fromLocalJSON{path: path}
}
