package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var dataDir string

func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dataDir = filepath.Join(home, ".dsu")
	return os.MkdirAll(dataDir, 0755)
}

func load(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

func save(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
