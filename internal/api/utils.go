package api

import "github.com/Masterminds/semver/v3"

func parseVersions(raw []string) ([]*semver.Version, error) {
	parsedVersions := make([]*semver.Version, len(raw))
	for _, version := range raw {
		v, err := semver.NewVersion(version)
		if err != nil {
			return nil, err
		}
		parsedVersions = append(parsedVersions, v)
	}
	return parsedVersions, nil
}
