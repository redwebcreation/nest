package pkg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"

	"github.com/redwebcreation/nest/util"
)

var (
	ErrInvalidStrategy   = fmt.Errorf("strategy must be either local or remote")
	ErrInvalidProvider   = fmt.Errorf("provider must be either github, gitlab or bitbucket")
	ErrInvalidRepository = fmt.Errorf("invalid repository name")
	ErrEmptyBranch       = fmt.Errorf("branch name cannot be empty")
)

var Config = &ConfigLocator{}

type ConfigLocatorConfig struct {
	Strategy   string
	Provider   string
	Repository string
	Branch     string
	Dir        string
	Commit     string
}

type ConfigLocator struct {
	ConfigLocatorConfig
	Git    *util.Repository
	config *Configuration
}

func (l *ConfigLocator) Retrieve() (*Configuration, error) {
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

func (l ConfigLocator) Read(path string) ([]byte, error) {
	if l.Dir != "" {
		path = strings.TrimSuffix(l.Dir, "/") + "/" + path
	}

	return l.Git.Read(path)
}

func (l ConfigLocator) GetRepositoryLocation() string {
	return fmt.Sprintf("git@%s.com:%s", l.Provider, l.Repository)
}

func (l ConfigLocator) cachePath() string {
	return "/tmp/" + base64.StdEncoding.EncodeToString([]byte(l.GetRepositoryLocation()))
}

func (l *ConfigLocator) UnmarshalJSON(data []byte) error {
	var lc ConfigLocatorConfig

	err := json.Unmarshal(data, &lc)
	if err != nil {
		return err
	}

	l.Strategy = lc.Strategy
	l.Provider = lc.Provider
	l.Repository = lc.Repository
	l.Dir = lc.Dir
	l.Branch = lc.Branch

	err = l.Validate()
	if err != nil {
		return err
	}

	repoPath := l.cachePath()
	var repo *util.Repository

	if _, err = os.Stat(repoPath); err != nil {
		repo, err = util.NewRepository(l.GetRepositoryLocation(), repoPath)
		if err != nil {
			return err
		}
	} else {
		repo, err = util.OpenRepository(repoPath)
		if err != nil {
			return err
		}
	}

	if l.Commit == "" {
		commit, err := repo.LatestCommit()
		if err != nil {
			return err
		}

		l.Commit = string(commit)

		err = repo.Checkout(l.Commit)
		if err != nil {
			return err
		}
	}

	l.Git = repo

	return nil
}

func (l ConfigLocator) Validate() error {
	if l.Strategy != "local" && l.Strategy != "remote" {
		return ErrInvalidStrategy
	}

	if l.Provider != "github" && l.Provider != "gitlab" && l.Provider != "bitbucket" {
		return ErrInvalidProvider
	}

	if l.Branch == "" {
		return ErrEmptyBranch
	}

	re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+(.git)?")
	if !re.MatchString(l.Repository) {
		return ErrInvalidRepository
	}

	return nil
}
