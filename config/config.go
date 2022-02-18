package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redwebcreation/nest/loggy"
	"gopkg.in/yaml.v2"
	"io/fs"
	"log"
	"os"
	"strings"
)

// Config contains nest's configuration
type Config struct {
	Provider   string `json:"provider"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Commit     string `json:"commit"`
	// StoreDir is the path where server configs are stored
	StoreDir string `json:"-"`
	// Path is the location of the config file
	Path   string `json:"-"`
	Logger *log.Logger
	Git    *Git
}

func (c *Config) StorePath() string {
	return c.StoreDir + "/" + c.Branch + "-" + strings.Replace(c.Repository, "/", "-", -1)
}

func (c *Config) RemoteURL() string {
	return fmt.Sprintf("git@%s.com:%s.git", c.Provider, c.Repository)
}

func (c *Config) Read(file string) ([]byte, error) {
	configPath := c.StorePath()

	_, err := os.Stat(configPath)

	if errors.Is(err, fs.ErrNotExist) {
		err = c.Clone()
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	c.log(loggy.DebugLevel, "reading serverConfig file", loggy.Fields{
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

	c.log(loggy.InfoLevel, "updating config", loggy.Fields{
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

	fmt.Printf("%+v\n", c)

	err := c.Git.Clone(c.RemoteURL(), c.StorePath(), c.Branch)

	if err != nil {
		return err
	}

	c.log(loggy.InfoLevel, "cloned config", loggy.Fields{
		"tag": "config.clone",
	})

	return nil
}

func (c *Config) log(level loggy.Level, message string, fields loggy.Fields) {
	fields["commit"] = c.Commit
	fields["branch"] = c.Branch
	fields["location"] = c.RemoteURL()

	c.Logger.Print(loggy.NewEvent(level, message, fields))
}

func (c *Config) Pull() error {
	_, err := c.Git.Pull(c.StorePath(), c.Branch)

	if err != nil {
		return err
	}

	c.log(loggy.InfoLevel, "pulled config", loggy.Fields{
		"tag": "config.pull",
	})

	return nil
}

// NewConfig creates a new config
// It isn't used while testing, make sure to reflect the changes you make here in the tests using the Config.
func NewConfig(configPath string, storePath string, log *log.Logger) (*Config, error) {
	contents, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("run `nest setup` to setup nest")
	}

	config := &Config{
		Path:     configPath,
		StoreDir: storePath,
		Logger:   log,
		Git: &Git{
			Logger: log,
		},
	}
	err = json.Unmarshal(contents, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
