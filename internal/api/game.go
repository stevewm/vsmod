package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Masterminds/semver/v3"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const GAME_API_BASE_URL = "https://api.vintagestory.at/"

type GameRelease struct {
	Version string `json:"version"`
}

type GameAPI struct {
	Client *http.Client
}

func NewGameAPI(client *http.Client) *GameAPI {
	return &GameAPI{
		Client: client,
	}
}

func (api *GameAPI) GetGameReleases(channel string) ([]*semver.Version, error) {
	if channel == "" {
		channel = "stable"
	}
	validChannels := []string{"stable", "unstable"}

	if !slices.Contains(validChannels, channel) {
		return nil, fmt.Errorf("invalid channel: %s. Valid channels are: %v", channel, validChannels)
	}

	resp, err := http.Get(GAME_API_BASE_URL + channel + ".json")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch game releases: %v", err)
	}

	defer func(b io.ReadCloser) {
		if b == nil {
			return
		}
		err = errors.Join(err, b.Close())
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var releases map[string]GameRelease
	if err := json.Unmarshal(body, &releases); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	raw := maps.Keys(releases)
	keys := make([]*semver.Version, len(raw))
	for i, r := range raw {
		v, err := semver.NewVersion(r)
		if err != nil {
			return nil, fmt.Errorf("failed to parse version: %v", err)
		}
		keys[i] = v
	}
	return keys, nil
}

func (api *GameAPI) GetLatestGameRelease(channel string) (*semver.Version, error) {
	releases, err := api.GetGameReleases(channel)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch game releases: %v", err)
	}

	return releases[0], nil
}
