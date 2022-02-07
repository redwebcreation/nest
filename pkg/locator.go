package pkg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/redwebcreation/nest/util"
)

var (
	ErrRepositoryNotFound    = fmt.Errorf("repository not found")
	ErrInvalidStrategy       = fmt.Errorf("strategy must be either local or remote")
	ErrInvalidProvider       = fmt.Errorf("provider must be either github, gitlab or bitbucket")
	ErrInvalidRepositoryName = fmt.Errorf("invalid repository name")
	ErrEmptyBranch           = fmt.Errorf("branch name cannot be empty")
)

var Config = &Locator{}

type Locator struct {
	Strategy   string
	Provider   string
	Repository string
	Branch     string
	Dir        string
	Commit     string
	repo       util.Repository
	config     *Configuration
}

func (l *Locator) Resolve() (*Configuration, error) {
	if l.config != nil {
		return l.config, nil
	}

	contents, err := l.Read("nest.yaml")
	if err != nil {
		return nil, err
	}

	var config Configuration

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (l Locator) Read(path string) ([]byte, error) {
	if l.Dir != "" {
		path = strings.TrimSuffix(l.Dir, "/") + "/" + path
	}

	repo, err := l.Repo()
	if err != nil {
		return nil, err
	}

	return repo.Read(path)
}

func (l Locator) GetRepositoryLocation() string {
	return fmt.Sprintf("repo@%s.com:%s", l.Provider, l.Repository)
}

func (l Locator) cachePath() string {
	return "/tmp/" + base64.StdEncoding.EncodeToString([]byte(l.GetRepositoryLocation()))
}

func (l *Locator) UnmarshalJSON(data []byte) error {
	type plain Locator
	var p plain

	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	l.Strategy = p.Strategy
	l.Provider = p.Provider
	l.Repository = p.Repository
	l.Dir = p.Dir
	l.Branch = p.Branch

	err = l.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (l Locator) Validate() error {
	if l.Strategy != "local" && l.Strategy != "remote" {
		return ErrInvalidStrategy
	}

	if l.Provider != "github" && l.Provider != "gitlab" && l.Provider != "bitbucket" {
		return ErrInvalidProvider
	}

	if l.Branch == "" {
		return ErrEmptyBranch
	}

	re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+(.repo)?")
	if !re.MatchString(l.Repository) {
		return ErrInvalidRepositoryName
	}

	return nil
}

func (l *Locator) Repo() (util.Repository, error) {
	if l.repo != nil {
		return l.repo, nil
	}

	repoPath := l.cachePath()
	var repo util.Repository

	if _, err := os.Stat(repoPath); err != nil {
		repo, err = util.NewRepository(l.GetRepositoryLocation(), repoPath)
		if err != nil {
			return nil, err
		}
	} else {
		repo, err = util.OpenRepository(repoPath)
		if err != nil {
			return nil, err
		}
	}

	if l.Branch != "" {
		err := repo.Checkout(l.Branch)
		if err != nil {
			return nil, err
		}
	}

	if l.Commit == "" {
		commit, err := repo.LatestCommit()
		if err != nil {
			return nil, err
		}

		l.Commit = string(commit)

		err = repo.Checkout(l.Commit)
		if err != nil {
			return nil, err
		}
	}

	l.repo = repo

	return repo, nil
}
