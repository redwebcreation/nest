package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrNoReleases      = fmt.Errorf("no releases found")
	ErrReleaseNotFound = fmt.Errorf("release not found")
)
var (
	GithubAPIURL = "https://api.github.com"
)

type GithubRepository string

type GithubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		State              string `json:"state"`
		BrowserDownloadUrl string `json:"browser_download_url"`
	}
}

func (g GithubRepository) newRequest(url string, data interface{}) error {
	response, err := http.Get(GithubAPIURL + "/" + url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &data)
}

func (g GithubRepository) Release(version string) (*GithubRelease, error) {
	var releases []GithubRelease
	err := g.newRequest(
		fmt.Sprintf("repos/%s/releases", g),
		&releases,
	)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, ErrNoReleases
	}

	if version == "" {
		return &releases[0], nil
	}

	for _, release := range releases {
		if release.TagName == version {
			return &release, nil
		}

	}

	return nil, ErrReleaseNotFound
}
