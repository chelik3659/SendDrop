package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Port     int    `json:"port"`
	IP       string `json:"ip"`
	Language string `json:"language"` // en, ru
	Theme    string `json:"theme"`    // light, dark
	ShareDir string `json:"share_dir"`
}

var (
	AppConfig *Config
	ConfigPath string
)

func Init() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	
	ConfigPath = filepath.Join(filepath.Dir(exePath), "config.json")
	
	// Загружаем или создаём дефолтный конфиг
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		AppConfig = DefaultConfig()
		return Save()
	}
	
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, &AppConfig)
}

func DefaultConfig() *Config {
	return &Config{
		Port:     8081,
		IP:       "0.0.0.0",
		Language: "en",
		Theme:    "dark",
		ShareDir: "sharing",
	}
}

func Save() error {
	data, err := json.MarshalIndent(AppConfig, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath, data, 0644)
}