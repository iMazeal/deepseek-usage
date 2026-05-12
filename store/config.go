package store

import "path/filepath"

func configPath() string { return filepath.Join(dataDir, "config.json") }

func SetConfig(key, value string) error {
	cfg := map[string]string{}
	load(configPath(), &cfg)
	cfg[key] = value
	return save(configPath(), cfg)
}

func GetConfig(key string) (string, error) {
	cfg := map[string]string{}
	if err := load(configPath(), &cfg); err != nil {
		return "", err
	}
	return cfg[key], nil
}
