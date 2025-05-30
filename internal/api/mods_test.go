package api

import (
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

var testClient = &http.Client{
	Timeout: 10 * time.Second,
}

var testMod = Mod{
	Id:      12345,
	Name:    "Test Mod",
	Desc:    "pls ignore",
	Author:  "stevewm",
	Created: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC).String(),
	Releases: ModReleaseSlice{
		{Version: *semver.MustParse("1.2")},
		{Version: *semver.MustParse("1.1")},
		{Version: *semver.MustParse("1.0")},
	},
}

var modAPI = NewModAPI(testClient)

func TestGetMod(t *testing.T) {
	mods := []string{"pspeltpatch", "3303"}
	for _, mod := range mods {
		modInfo, err := modAPI.GetMod(mod)
		assert.NoError(t, err, "Expected no error when fetching mod info")
		assert.NotNil(t, modInfo, "Expected mod info to be not nil")
		assert.True(t, semver.MustParse("0.0.2").Equal(&modInfo.Releases[0].Version), "Expected mod version to match")
	}
}

func TestModLatestConstraintRelease(t *testing.T) {
	constraint, err := semver.NewConstraint("~1.0")
	assert.NoError(t, err, "Expected no error creating constraint")

	latestRelease, err := testMod.Release(*constraint)
	assert.NoError(t, err, "Expected no error when getting latest release with constraint")
	assert.True(t, constraint.Check(&latestRelease.Version), "Expected latest release version to be 1.2")

	constraintNoMatch, err := semver.NewConstraint(">=2.0")
	assert.NoError(t, err, "Expected no error creating constraint")
	release, err := testMod.Release(*constraintNoMatch)
	assert.Error(t, err, "Expected error when no releases match the constraint")
	assert.Nil(t, release, "Expected no release when constraint does not match any release")
}

func TestModReleaseCompatibleWithGameVersion(t *testing.T) {
	testRelease := ModRelease{
		Tags: []semver.Version{*semver.MustParse("1.0"), *semver.MustParse("1.1")},
	}
	type test struct {
		gameVersion string
		compatible  bool
	}
	tc := []test{
		{gameVersion: "1.0", compatible: true},
		{gameVersion: "1.1", compatible: true},
		{gameVersion: "1.2", compatible: false},
		{gameVersion: "2.0", compatible: false},
		// v-prefix
		{gameVersion: "v1.0", compatible: true},
		{gameVersion: "v1.1", compatible: true},
		{gameVersion: "v1.2", compatible: false},
		{gameVersion: "v2.0", compatible: false},
	}

	for _, testCase := range tc {
		constraint, err := semver.NewConstraint(testCase.gameVersion)
		assert.NoError(t, err, "Expected no error creating constraint for game version %s", testCase.gameVersion)
		isCompatible := testRelease.CompatibleWith(*constraint)
		assert.Equal(t, testCase.compatible, isCompatible, "Expected mod release compatibility with game version %s to be %v", testCase.gameVersion, testCase.compatible)
	}
}

func TestModLatestRelease(t *testing.T) {
	latestRelease, err := testMod.LatestRelease()
	assert.NoError(t, err, "Expected no error when getting latest release")
	assert.True(t, semver.MustParse("1.2").Equal(&latestRelease.Version), "Expected latest release version to be 1.2")
}

func TestModReleaseSorting(t *testing.T) {
	sortedReleases := ModReleaseSlice{
		{Version: *semver.MustParse("1.0")},
		{Version: *semver.MustParse("1.2")},
		{Version: *semver.MustParse("1.1")},
	}

	sort.Sort(sortedReleases)

	assert.Len(t, sortedReleases, 3, "Expected 3 releases after sorting")
	assert.True(t, sortedReleases[0].Version.Equal(semver.MustParse("1.2")), "Expected first release to be 1.2")
	assert.True(t, sortedReleases[0].Version.Equal(semver.MustParse("1.2")), "Expected second release to be 1.2")
	assert.True(t, sortedReleases[0].Version.Equal(semver.MustParse("1.2")), "Expected third release to be 1.2")
}

func TestDownloadURL(t *testing.T) {
	modRelease := ModRelease{
		URL: "https://example.com/mod file.zip",
	}

	assert.Equal(t, "https://example.com/mod%20file.zip", modRelease.DownloadURL(), "Expected download URL to be properly encoded")
	assert.False(t, strings.Contains(modRelease.DownloadURL(), " "), "Expected download URL to not contain spaces")
}
