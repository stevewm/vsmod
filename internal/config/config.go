package config

import (
	"path/filepath"
)

type ConfigFile struct {
	ModsDir     string                  `yaml:"mods_dir"`
	GameVersion string                  `yaml:"game_version"`
	Mods        []ConfigFileMod         `yaml:"mods"`
	Hooks       map[string]CommandHooks `yaml:"hooks"`
}

type ConfigFileMod struct {
	ID          string `yaml:"id"`
	Version     string `yaml:"version"`
	CompatCheck bool   `yaml:"compatibility_check"`
}

func (c *ConfigFile) Dir() string {
	absPath, _ := filepath.Abs(c.ModsDir)
	return absPath
}
