package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testClient = &http.Client{
	Timeout: 10 * time.Second,
}
var modAPI = NewModAPI(testClient)

func TestGetMod(t *testing.T) {
	mods := []string{"pspeltpatch", "3303"}
	for _, mod := range mods {
		modInfo, err := modAPI.GetMod(mod)
		assert.NoError(t, err, "Expected no error when fetching mod info")
		assert.NotNil(t, modInfo, "Expected mod info to be not nil")
		assert.Equal(t, "0.0.2", modInfo.Releases[0].Version, "Expected mod version to match") // Fragile
	}
}

func TestModReleaseCompatibleWithGameVersion(t *testing.T) {
	modRelease := ModRelease{
		Tags: []string{"v1.0", "v1.1"},
	}

	assert.True(t, modRelease.CompatibleWith("1.0"), "Expected mod release to be compatible with game version")
	assert.False(t, modRelease.CompatibleWith("2.0"), "Expected mod release to not be compatible with game version")
	assert.False(t, modRelease.CompatibleWith(""), "Expected empty version to not be compatible")
}

func TestDownloadURL(t *testing.T) {
	modRelease := ModRelease{
		URL: "https://example.com/mod file.zip",
	}

	assert.Equal(t, "https://example.com/mod%20file.zip", modRelease.DownloadURL(), "Expected download URL to be properly encoded")
	assert.NotEqual(t, "https://example.com/mod file.zip", modRelease.DownloadURL(), "Expected download URL to be properly encoded")
}
