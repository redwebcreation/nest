package util

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGithubRepository_Release(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	GithubAPIURL = srv.URL

	release, err := GithubRepository("redwebcreation/nest").Release("")
	if err != nil {
		t.Fatal(err)
	}

	if release.TagName != "v0.1.0" {
		t.Errorf("Expected v0.1.0, got %s", release.TagName)
	}

	release, err = GithubRepository("redwebcreation/nest").Release("v0.1.0")
	if err != nil {
		t.Fatal(err)
	}

	if release.TagName != "v0.1.0" {
		t.Errorf("Expected v0.1.0, got %s", release.TagName)
	}

	release, err = GithubRepository("redwebcreation/nest").Release("v0.0.2")
	if err != ErrReleaseNotFound {
		t.Errorf("Expected ErrReleaseNotFound, got %s", err)
	}

	release, err = GithubRepository("redwebcreation/nest").Release("v1.0.0")
	if err != nil {
		t.Fatal(err)
	}

	if release.TagName != "v1.0.0" {
		t.Errorf("Expected v1.0.0, got %s", release.TagName)
	}

	if len(release.Assets) != 1 {
		t.Errorf("Expected 1 asset, got %d", len(release.Assets))
	}

	if release.Assets[0].State != "uploaded" {
		t.Errorf("Expected uploaded, got %s", release.Assets[0].State)
	}

	if release.Assets[0].BrowserDownloadUrl != "example.com/download" {
		t.Errorf("Expected example.com/download, got %s", release.Assets[0].BrowserDownloadUrl)
	}
}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()

	handler.HandleFunc("/repos/redwebcreation/nest/releases", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"tag_name": "v0.1.0"}, {"tag_name": "v1.0.0", "assets": [{"state": "uploaded", "browser_download_url": "example.com/download"}]}]`))
	})

	srv := httptest.NewServer(handler)

	return srv
}
