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

var testMod = Mod{
	Releases: []ModRelease{
		{Version: "1.2"},
		{Version: "1.1"},
		{Version: "1.0"},
	},
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

func TestModLatestRelease(t *testing.T) {
	latestRelease, err := testMod.LatestRelease()
	assert.NoError(t, err, "Expected no error when getting latest release")
	assert.Equal(t, "1.2", latestRelease.Version, "Expected latest release version to be 1.2")
}

func TestGetModRelease(t *testing.T) {
	release, err := testMod.Release("1.1")
	assert.NoError(t, err, "Expected no error when getting specific release")
	assert.NotNil(t, release, "Expected release to be not nil")
	assert.Equal(t, "1.1", release.Version, "Expected release version to be 1.1")

	release, err = testMod.Release("2.0")
	assert.Error(t, err, "Expected error when getting non-existent release")
	assert.Nil(t, release, "Expected release to be nil for non-existent version")
}

func TestDownloadURL(t *testing.T) {
	modRelease := ModRelease{
		URL: "https://example.com/mod file.zip",
	}

	assert.Equal(t, "https://example.com/mod%20file.zip", modRelease.DownloadURL(), "Expected download URL to be properly encoded")
	assert.NotEqual(t, "https://example.com/mod file.zip", modRelease.DownloadURL(), "Expected download URL to be properly encoded")
}
