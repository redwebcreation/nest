package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"gopkg.in/yaml.v2"
	"io/fs"
	"log"
	"os"
	"strings"
)

type Config struct {
	Provider   string `json:"provider"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Commit     string `json:"commit"`
	// Store is the path where configs are stored
	Store  string `json:"-"`
	Path   string `json:"-"`
	logger *log.Logger
	Git    *git
}

func (c *Config) StorePath() string {
	return c.Store + "/" + c.Branch + "-" + strings.Replace(c.Repository, "/", "-", -1)
}

func (c *Config) RemoteURL() string {
	return fmt.Sprintf("git@%s.com:%s.git", c.Provider, c.Repository)
}

func (c *Config) Read(file string) ([]byte, error) {
	configPath := c.StorePath()

	if _, err := os.Stat(configPath); errors.Is(err, fs.ErrNotExist) {
		err = c.Clone()
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	c.log(global.LevelDebug, "reading serverConfig file", global.Fields{
		"tag":  "ServerConfig.read",
		"file": file,
	})

	return c.Git.ReadFile(configPath, c.Commit, file)
}

func (c *Config) ServerConfig() (*ServerConfig, error) {
	configFile := "nest.yaml"
	if c.Git.Exists(c.StorePath(), "nest.yml", c.Commit) {
		configFile = "nest.yml"
	}

	contents, err := c.Read(configFile)
	if err != nil {
		return nil, err
	}

	config := &ServerConfig{}
	err = yaml.Unmarshal(contents, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Save() error {
	contents, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(c.Path, contents, 0600)
	if err != nil {
		return err
	}

	c.log(global.LevelInfo, "updating config", global.Fields{
		"tag": "config.update",
	})

	return nil
}

func (c *Config) LoadCommit(commit string) error {
	c.Commit = commit

	return c.Save()
}

func (c *Config) Clone() error {
	_ = os.RemoveAll(c.StorePath())

	err := c.Git.Clone(c.RemoteURL(), c.StorePath(), c.Branch)

	if err != nil {
		return err
	}

	c.log(global.LevelInfo, "cloned config", global.Fields{
		"tag": "config.clone",
	})

	return nil
}

func (c *Config) log(level global.Level, message string, fields global.Fields) {
	fields["commit"] = c.Commit
	fields["branch"] = c.Branch
	fields["location"] = c.RemoteURL()

	c.logger.Print(global.NewEvent(level, message, fields))
}

func (c *Config) Pull() error {
	_, err := c.Git.Pull(c.StorePath(), c.Branch)

	if err != nil {
		return err
	}

	c.log(global.LevelInfo, "pulled config", global.Fields{
		"tag": "config.pull",
	})

	return nil
}

func NewConfig(configPath string, storePath string, logger *log.Logger) (*Config, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	contents, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("run `nest setup` to setup nest")
	}

	config := &Config{
		Path:   configPath,
		Store:  storePath,
		logger: logger,
		Git: &git{
			logger: logger,
		},
	}
	err = json.Unmarshal(contents, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
