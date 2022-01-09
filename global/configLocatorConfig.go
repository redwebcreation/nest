package global

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/mitchellh/go-homedir"
)

var LocatorConfigFile string
var IsConfigLocatorConfigured bool
var ConfigLocatorConfig *LocatorConfig

type locatorConfig struct {
	Strategy    string `json:"strategy,omitempty"`
	ProviderURL string `json:"provider_url,omitempty"`
	Repository  string `json:"repository,omitempty"`
	Cache       string `json:"cache,omitempty"`
}

type LocatorConfig struct {
	locatorConfig
	Head     *object.Commit `json:"head,omitempty"`
	Worktree *git.Worktree  `json:"worktree,omitempty"`
}

func (r LocatorConfig) SaveLocally() error {
	contents, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return os.WriteFile(LocatorConfigFile, contents, 0644)
}

func (r LocatorConfig) Read(path string) ([]byte, error) {
	f, err := r.Worktree.Filesystem.Open(path)
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

func (r LocatorConfig) Exists(path string) bool {
	_, err := r.Worktree.Filesystem.Stat(path)
	return err != nil
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	LocatorConfigFile = home + "/.nest.json"

	_, err = os.Stat(LocatorConfigFile)
	IsConfigLocatorConfigured = err == nil

	var f LocatorConfig
	if IsConfigLocatorConfigured {
		contents, err := os.ReadFile(LocatorConfigFile)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(contents, &f)
		if err != nil {
			panic(err)
		}
	}
	ConfigLocatorConfig = &f
}

func (r *LocatorConfig) UnmarshalJSON(data []byte) error {
	var rc locatorConfig
	err := json.Unmarshal(data, &rc)
	if err != nil {
		return err
	}

	r.Strategy = rc.Strategy
	r.ProviderURL = rc.ProviderURL
	r.Repository = rc.Repository
	r.Cache = rc.Cache

	if r.Cache == "" {
		r.Cache = "/tmp"
	}

	repositoryPath := strings.TrimSuffix(r.Cache, "/") + "/" + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s-%s", r.ProviderURL, r.Repository)))

	var repo *git.Repository

	if _, existsErr := os.Stat(repositoryPath); existsErr != nil {
		repo, err = git.PlainClone(repositoryPath, false, &git.CloneOptions{
			URL: "git@github.com:" + r.Repository,
		})
	} else {
		repo, err = git.PlainOpen(repositoryPath)
	}
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	ref, _ := repo.Head()
	commit, _ := repo.CommitObject(ref.Hash())

	_ = worktree.Checkout(&git.CheckoutOptions{
		Hash: commit.Hash,
	})

	r.Head = commit
	r.Worktree = worktree
	return nil
}
