package config

import (
	"path/filepath"

	"github.com/Masterminds/semver/v3"
)

type ConfigFile struct {
	ModsDir     string                  `yaml:"mods_dir"`
	GameVersion semver.Constraints      `yaml:"game_version"`
	Mods        []ConfigFileMod         `yaml:"mods"`
	Hooks       map[string]CommandHooks `yaml:"hooks"`
}

type ConfigFileMod struct {
	ID          string             `yaml:"id"`
	Version     semver.Constraints `yaml:"version"`
	CompatCheck bool               `yaml:"compatibility_check"`
}

func (c *ConfigFile) Dir() string {
	absPath, _ := filepath.Abs(c.ModsDir)
	return absPath
}
