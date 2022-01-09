package common

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/me/nest/global"
	"io"
	"os"
	"strings"
)

var ConfigReader *configReader

type LocatorConfig struct {
	Strategy    string
	ProviderURL string
	Repository  string
	Cache       string
}

type configReader struct {
	LocatorConfig
	Head     *object.Commit
	Worktree *git.Worktree `json:"-"`
}

func (c configReader) WriteOnDisk() error {
	contents, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(global.ConfigLocatorConfigFile, contents, 0644)
}

func (c configReader) Read(path string) ([]byte, error) {
	f, err := c.Worktree.Filesystem.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func LoadConfigReader() (*configReader, error) {
	var cr configReader

	contents, err := os.ReadFile(global.ConfigLocatorConfigFile)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(contents, &cr); err != nil && err.Error() == "unknown error: remote: " {
		return nil, fmt.Errorf("the repository %s does not exists", cr.ProviderURL+cr.Repository)
	} else {
		return &cr, err
	}
}

func (c configReader) getCacheKey() string {
	return base64.StdEncoding.EncodeToString([]byte(c.ProviderURL + c.Repository))
}

func (c *configReader) UnmarshalJSON(data []byte) error {
	var lc LocatorConfig
	err := json.Unmarshal(data, &lc)
	if err != nil {
		return err
	}

	c.Strategy = lc.Strategy
	c.ProviderURL = lc.ProviderURL
	c.Repository = lc.Repository
	c.Cache = lc.Cache

	if c.Cache == "" {
		c.Cache = "/tmp"
	} else {
		c.Cache = strings.TrimSuffix(c.Cache, "/")
	}

	repositoryPath := c.Cache + "/" + c.getCacheKey()

	var repo *git.Repository

	if _, e := os.Stat(repositoryPath); e != nil {
		repo, err = git.PlainClone(repositoryPath, false, &git.CloneOptions{
			URL: c.ProviderURL + c.Repository + ".git",
		})
		if err != nil {
			return err
		}
	} else {
		repo, err = git.PlainOpen(repositoryPath)
		if err != nil {
			return err
		}
	}

	worktree, _ := repo.Worktree()
	ref, _ := repo.Head()
	commit, _ := repo.CommitObject(ref.Hash())

	_ = worktree.Checkout(&git.CheckoutOptions{
		Hash: commit.Hash,
	})

	c.Head = commit
	c.Worktree = worktree
	return nil
}
