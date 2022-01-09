package common

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/util"
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

	repo := &util.Repository{
		Path: c.Cache + "/" + c.getCacheKey(),
	}
	if _, e := os.Stat(repo.Path); e != nil {
		err = repo.Clone(c.ProviderURL + c.Repository)
		if err != nil {
			return err
		}
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
