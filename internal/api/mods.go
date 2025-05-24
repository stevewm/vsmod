package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"
)

const MOD_API_BASE_URL = "https://mods.vintagestory.at/api/"

// for some reason the moddb returns a mod's info as the value of "mod"
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
	defer resp.Body.Close()

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

type ModWrapper struct {
	Mod Mod `json:"mod"`
}

type Mod struct {
	Id       int          `json:"modid"`
	Name     string       `json:"name"`
	Desc     string       `json:"text"`
	Author   string       `json:"author"`
	Homepage string       `json:"homepageurl"`
	Created  string       `json:"created"`
	Releases []ModRelease `json:"releases"`
}

type ModRelease struct {
	Created  string   `json:"created"`
	FileName string   `json:"filename"`
	Tags     []string `json:"tags"`
	URL      string   `json:"mainfile"`
	Version  string   `json:"modversion"`
}

func (m *Mod) LatestRelease() (*ModRelease, error) {
	// fixme: dont rely on the moddb ordering releases properly
	if len(m.Releases) == 0 {
		return nil, fmt.Errorf("mod has no releases")
	}
	return &m.Releases[0], nil
}

func (m *Mod) Release(version string) (*ModRelease, error) {
	if version == "" {
		latestRelease, err := m.LatestRelease()

		if err != nil {
			return nil, err
		}
		return latestRelease, nil
	}

	for _, release := range m.Releases {
		if release.Version == version {
			return &release, nil
		}
	}
	return nil, fmt.Errorf("version not found")
}

func (m *ModRelease) DownloadURL() string {
	url := strings.ReplaceAll(m.URL, " ", "%20")
	return url
}

func (m *ModRelease) CompatibleWith(gameVersion string) bool {
	// todo: convert to semver
	if !strings.HasPrefix(gameVersion, "v") {
		gameVersion = "v" + gameVersion
	}
	return slices.Contains(m.Tags, gameVersion)
}
