package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	CacheDir    string   `json:"cache_dir"`
	EnabledApps []string `json:"enabled_apps"`
	SwwwEnabled bool     `json:"swww_enabled"`
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cacheDir := filepath.Join(homeDir, ".cache", "gothemr")
	configPath := filepath.Join(homeDir, ".config", "gothemr", "config.json")

	cfg := &Config{
		CacheDir:    cacheDir,
		EnabledApps: getDefaultApps(),
		SwwwEnabled: true,
	}

	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	if err := os.MkdirAll(cfg.CacheDir, 0755); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getDefaultApps() []string {
	return []string{
		"waybar",
		"rofi",
		"hyprland",
		"gtk",
		"terminal",
	}
}

func generateCacheKey(inputPath string, colorCount int) string {
	return fmt.Sprintf("%s_%d", filepath.Base(inputPath), colorCount)
}

func shouldRegenerate(themePath, inputPath string) bool {
	themeInfo, err := os.Stat(themePath)
	if err != nil {
		return true
	}

	inputInfo, err := os.Stat(inputPath)
	if err != nil {
		return true
	}

	return inputInfo.ModTime().After(themeInfo.ModTime())
}
