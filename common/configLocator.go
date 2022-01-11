package common

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/util"
)

var ConfigReader *configReader

type LocatorConfig struct {
	Strategy   string
	Provider   string
	Repository string
}

type configReader struct {
	LocatorConfig
	LatestCommit string
	Git          *util.Repository
}

func (c configReader) WriteOnDisk() error {
	contents, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(global.ConfigLocatorConfigFile, contents, 0644)
}

func (c configReader) Read(path string) ([]byte, error) {
	data, err := c.Git.ReadFile(path)

	return []byte(data), err
}

func (c configReader) GetRepositoryLocation() string {
	return fmt.Sprintf("git@%s.com:%s", c.Provider, c.Repository)
}

func LoadConfigReader() (*configReader, error) {
	var cr configReader

	contents, err := os.ReadFile(global.ConfigLocatorConfigFile)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(contents, &cr); err != nil && err.Error() == "unknown error: remote: " {
		return nil, fmt.Errorf("the repository %s does not exists", cr.GetRepositoryLocation())
	} else {
		return &cr, err
	}
}

func (c configReader) getCacheKey() string {
	return base64.StdEncoding.EncodeToString([]byte(c.GetRepositoryLocation()))
}

func (c *configReader) UnmarshalJSON(data []byte) error {
	var lc LocatorConfig
	err := json.Unmarshal(data, &lc)
	if err != nil {
		return err
	}

	c.Strategy = lc.Strategy
	c.Provider = lc.Provider
	c.Repository = lc.Repository

	err = c.Validate()
	if err != nil {
		return err
	}

	repo := &util.Repository{
		Path: "/tmp/" + c.getCacheKey(),
	}
	if _, err := os.Stat(repo.Path); errors.Is(err, os.ErrNotExist) {
		err = repo.Clone(c.GetRepositoryLocation())
		if err != nil {
			return err
		}
	} else {
		return err
	}

	commit, err := repo.LatestCommit()
	if err != nil {
		return err
	}

	err = repo.Checkout(commit)
	if err != nil {
		return err
	}

	c.LatestCommit = commit
	c.Git = repo

	return nil
}

var (
	ErrInvalidStrategy   = fmt.Errorf("strategy must be either local or remote")
	ErrInvalidProvider   = fmt.Errorf("provider must be either github, gitlab or bitbucket")
	ErrInvalidRepository = fmt.Errorf("invalid repository name")
)

func (c configReader) Validate() error {
	if c.Strategy != "local" && c.Strategy != "remote" {
		return ErrInvalidStrategy
	}

	if c.Provider != "github" && c.Provider != "gitlab" && c.Provider != "bitbucket" {
		return ErrInvalidProvider
	}

	re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+(.git)?")
	if !re.MatchString(c.Repository) {
		return ErrInvalidRepository
	}

	return nil
}

func NewConfigReader() *configReader {
	return &configReader{}
}
