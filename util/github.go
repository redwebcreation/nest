package util

import (
	"encoding/json"
	"io"
	"net/http"
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
	response, err := http.Get("https://api.github.com/" + url)
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

func (g GithubRepository) String() string {
	return string(g)
}

func (g GithubRepository) Releases(version string) ([]GithubRelease, error) {
	var releases []GithubRelease
	err := g.newRequest("repos/"+g.String()+"/releases", &releases)
	if err != nil {
		return nil, err
	}
	var filtered []GithubRelease

	for _, release := range releases {
		if release.TagName != version && version != "" {
			continue
		}

		filtered = append(filtered, release)
	}

	return filtered, nil
}
