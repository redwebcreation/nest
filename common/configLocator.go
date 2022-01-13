package common

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/redwebcreation/nest/util"
)

var (
	ErrInvalidStrategy   = fmt.Errorf("strategy must be either local or remote")
	ErrInvalidProvider   = fmt.Errorf("provider must be either github, gitlab or bitbucket")
	ErrInvalidRepository = fmt.Errorf("invalid repository name")
)

var ConfigLocator = &LocatorConfig{}

type LocatorConfig struct {
	Strategy   string
	Provider   string
	Repository string
	Branch     string
	Dir        string
	Commit     string
	Git        *util.Repository
}

func (l LocatorConfig) Read(path string) ([]byte, error) {
	if l.Dir != "" {
		path = strings.TrimSuffix(l.Dir, "/") + "/" + path
	}

	data, err := l.Git.Read(path)

	return []byte(data), err
}

func (l LocatorConfig) GetRepositoryLocation() string {
	return fmt.Sprintf("git@%s.com:%s", l.Provider, l.Repository)
}

func (l LocatorConfig) cachePath() string {
	return "/tmp/" + base64.StdEncoding.EncodeToString([]byte(l.GetRepositoryLocation()))
}

func (l *LocatorConfig) UnmarshalJSON(data []byte) error {
	var lc struct {
		Strategy   string
		Provider   string
		Repository string
		Branch     string
		Commit     string
		Dir        string
	}

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

		l.Commit = commit

		err = repo.Checkout(commit)
		if err != nil {
			return err
		}
	}

	l.Git = repo

	return nil
}

func (l LocatorConfig) Validate() error {
	if l.Strategy != "local" && l.Strategy != "remote" {
		return ErrInvalidStrategy
	}

	if l.Provider != "github" && l.Provider != "gitlab" && l.Provider != "bitbucket" {
		return ErrInvalidProvider
	}

	if l.Branch == "" {
		l.Branch = "main"
	}

	re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+(.git)?")
	if !re.MatchString(l.Repository) {
		return ErrInvalidRepository
	}

	return nil
}
