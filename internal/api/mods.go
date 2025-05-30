package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

const MOD_API_BASE_URL = "https://mods.vintagestory.at/api/"

type ModAPI struct {
	Client *http.Client
}

// Creates a new Vintage Story Mod API client
func NewModAPI(client *http.Client) *ModAPI {
	return &ModAPI{
		Client: client,
	}
}

// Gets details of a mod from the API
func (api *ModAPI) GetMod(modID string) (*Mod, error) {
	resp, err := api.Client.Get(MOD_API_BASE_URL + "mod/" + modID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mod info: %s: %w", modID, err)
	}
	defer func(b io.ReadCloser) {
		if b == nil {
			return
		}
		err = errors.Join(err, b.Close())
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body:  %s: %w", modID, err)
	}

	var wrapper ModWrapper
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON:  %s: %w", modID, err)
	}
	return &wrapper.Mod, nil
}

// the moddb returns a mod's info as the value of `mod` in the JSON response
type ModWrapper struct {
	Mod Mod `json:"mod"`
}

type Mod struct {
	Id       int             `json:"modid"`
	Name     string          `json:"name"`
	Desc     string          `json:"text"`
	Author   string          `json:"author"`
	Homepage string          `json:"homepageurl"`
	Created  string          `json:"created"`
	Releases ModReleaseSlice `json:"releases"`
}

type ModRelease struct {
	Created  string           `json:"created"`
	FileName string           `json:"filename"`
	Tags     []semver.Version `json:"tags"`
	URL      string           `json:"mainfile"`
	Version  semver.Version   `json:"modversion"`
}

type ModReleaseSlice []ModRelease

func (s ModReleaseSlice) Len() int {
	return len(s)
}

func (s ModReleaseSlice) Less(i, j int) bool {
	// descending order (latest release first)
	return s[j].Version.LessThan(&s[i].Version)
}

func (s ModReleaseSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// LatestRelease returns the latest release of the mod or an error if the mod has no releases.
func (m *Mod) LatestRelease() (*ModRelease, error) {
	if len(m.Releases) == 0 {
		return nil, fmt.Errorf("mod has no releases")
	}
	sort.Sort(m.Releases)
	return &m.Releases[0], nil
}

// Release finds the latest release of a mod that matches the given version constraint.
// It returns the first release that satisfies the constraint, or an error if no such release exists.
func (m *Mod) Release(modConstraint semver.Constraints) (*ModRelease, error) {
	log.Debugf("Finding latest release for mod %s with constraint %s", m.Name, modConstraint.String())

	if len(m.Releases) == 0 {
		return nil, fmt.Errorf("mod %s has no releases: %s", m.Name, m.Releases)
	}

	log.Debugf("Mod %s has %d releases\n", m.Name, len(m.Releases))

	validReleases := lo.Filter(m.Releases, func(release ModRelease, _ int) bool {
		log.Debugf("Checking release %s against constraint %s", release.Version.String(), modConstraint.String())
		return modConstraint.Check(&release.Version)
	})

	sort.Sort(m.Releases)

	log.Debugf("Found %d valid releases for mod %s that match constraint %s", len(validReleases), m.Name, modConstraint.String())

	if len(validReleases) > 0 {
		log.Debugf("valid releases: %v", validReleases)
		return &validReleases[0], nil
	}
	return nil, fmt.Errorf("no releases found for mod %s that match constraint %s", m.Name, modConstraint.String())
}

// DownloadURL returns the download URL for the mod release, with spaces replaced by %20.
func (m *ModRelease) DownloadURL() string {
	url := strings.ReplaceAll(m.URL, " ", "%20")
	return url
}

// CompatibleWith checks if the mod release is compatible with the given game version constraint.
func (m *ModRelease) CompatibleWith(gameVersion semver.Constraints) bool {
	log.Debugf("Checking mod release %s against game version %s", m.Version.String(), gameVersion.String())
	for _, tag := range m.Tags {
		ok := gameVersion.Check(&tag)
		log.Debugf("Checking tag %s against game version %s: %v", tag.String(), gameVersion.String(), ok)
		if ok {
			return true
		}
	}
	log.Warnf("mod release %s is not compatible with game version %s", m.Version.String(), gameVersion.String())
	return false
}
