package config

import (
	"path/filepath"
	"testing"
)

func TestAbsoluteModsDir(t *testing.T) {
	config := &ConfigFile{
		ModsDir: "mods",
	}

	absPath := config.Dir()

	expectedPath, _ := filepath.Abs("mods")
	if absPath != expectedPath {
		t.Errorf("Expected %s, got %s", expectedPath, absPath)
	}
}
