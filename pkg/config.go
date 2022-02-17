package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"gopkg.in/yaml.v2"
	"io/fs"
	"os"
	"strings"
)

type Config struct {
	Provider   string `json:"provider"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Commit     string `json:"commit"`
}

func (c *Config) Path() string {
	return global.ServerConfigsStore() + "/" + c.Branch + "-" + strings.Replace(c.Repository, "/", "-", -1)
}

func (c *Config) RemoteURL() string {
	return fmt.Sprintf("git@%s.com:%s.git", c.Provider, c.Repository)
}

func (c *Config) Read(file string) ([]byte, error) {
	configPath := c.Path()

	if _, err := os.Stat(configPath); errors.Is(err, fs.ErrNotExist) {
		err = c.Clone()
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	c.log(global.LevelDebug, "reading serverConfig file", global.Fields{
		"tag":  "ServerConfiguration.read",
		"file": file,
	})

	return Git.ReadFile(configPath, c.Commit, file)
}

func (c *Config) GetServerConfiguration() (*ServerConfiguration, error) {
	configFile := "nest.yaml"
	if Git.Exists(c.Path(), "nest.yml", c.Commit) {
		configFile = "nest.yml"
	}

	contents, err := c.Read(configFile)
	if err != nil {
		return nil, err
	}

	config := &ServerConfiguration{}
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
	err = os.WriteFile(global.ConfigFile(), contents, 0600)
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
	_ = os.RemoveAll(c.Path())

	err := Git.Clone(c.RemoteURL(), c.Path(), c.Branch)

	if err != nil {
		return err
	}

	c.log(global.LevelInfo, "cloned config", global.Fields{
		"tag": "config.clone",
	})

	return nil
}

func (c *Config) log(level global.Level, message string, fields global.Fields) {
	fields["Commit"] = c.Commit
	fields["branch"] = c.Branch
	fields["location"] = c.RemoteURL()

	global.LogI(level, message, fields)
}

func (c *Config) Pull() error {
	_, err := Git.Pull(c.Path(), c.Branch)

	if err != nil {
		return err
	}

	c.log(global.LevelInfo, "pulled config", global.Fields{
		"tag": "config.pull",
	})

	return nil
}

func NewConfig() (*Config, error) {
	contents, err := os.ReadFile(global.ConfigFile())
	if err != nil {
		return nil, fmt.Errorf("run `nest setup` to setup nest")
	}

	config := &Config{}
	err = json.Unmarshal(contents, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
