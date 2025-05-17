package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetGameReleases(t *testing.T) {
	var testClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	var gameAPI = NewGameAPI(testClient)
	releases, err := gameAPI.GetGameReleases("stable")

	assert.NoError(t, err, "Expected no error when fetching game releases")
	assert.NotEmpty(t, releases)
}
