package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const configDir = ".apialerts"
const configFile = "config.json"

// configDirOverride allows overriding the config directory for testing.
var configDirOverride string

type CLIConfig struct {
	APIKey    string `json:"api_key"`
	ServerURL string `json:"server_url,omitempty"`
}

func configPath() (string, error) {
	if configDirOverride != "" {
		return filepath.Join(configDirOverride, configFile), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, configFile), nil
}

func Load() (*CLIConfig, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &CLIConfig{}, nil
		}
		return nil, err
	}

	var cfg CLIConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(cfg *CLIConfig) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func GetAPIKey() (string, error) {
	cfg, err := Load()
	if err != nil {
		return "", err
	}
	if cfg.APIKey == "" {
		return "", errors.New("no API key configured — run: apialerts init")
	}
	return cfg.APIKey, nil
}
