package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

var Locator = &locator{}

type locator struct {
	Provider   string
	Repository string
	Branch     string
	Commit     string
}

func (l locator) ConfigPath() string {
	return global.GetConfigStoreDir() + "/" + strings.Replace(l.Repository, "/", "-", -1)
}

func (l locator) RemoteURL() string {
	return fmt.Sprintf("git@%s.com:%s.git", l.Provider, l.Repository)
}

func (l locator) Read(file string) ([]byte, error) {
	configPath := l.ConfigPath()

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		err = l.CloneConfig()
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	l.log(global.LevelDebug, "reading config file", global.Fields{
		"tag":  "locator.read",
		"file": file,
	})

	return Git.ReadFile(configPath, l.Commit, file)
}

func (l locator) Resolve() (*Configuration, error) {
	err := l.Load()
	if err != nil {
		return nil, err
	}

	configFile := "nest.yaml"
	if Git.Exists(l.ConfigPath(), "nest.yml", l.Commit) {
		configFile = "nest.yml"
	}

	contents, err := l.Read(configFile)
	if err != nil {
		return nil, err
	}

	config := &Configuration{}
	err = yaml.Unmarshal(contents, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (l *locator) Load() error {
	contents, err := os.ReadFile(global.GetLocatorConfigFile())
	if err != nil {
		return err
	}

	var p locator
	err = json.Unmarshal(contents, &p)
	if err != nil {
		return err
	}

	*l = p

	if l.Commit == "" {
		return fmt.Errorf("commit is empty, run `nest setup` to set it")
	}

	return nil
}

func (l *locator) LoadCommit(commit string) error {
	l.Commit = commit
	err := l.Save()
	if err != nil {
		return err
	}

	return l.Load()
}

func (l *locator) Save() error {
	contents, err := json.Marshal(l)
	if err != nil {
		return err
	}
	err = os.WriteFile(global.GetLocatorConfigFile(), contents, 0600)
	if err != nil {
		return err
	}

	l.log(global.LevelInfo, "updating locator config", global.Fields{
		"tag": "locator.update",
	})

	return nil
}

func (l locator) CloneConfig() error {
	_ = os.RemoveAll(l.ConfigPath())

	err := Git.Clone(l.RemoteURL(), l.ConfigPath())

	if err != nil {
		return err
	}

	l.log(global.LevelInfo, "cloned config", global.Fields{
		"tag": "locator.clone",
	})

	return nil
}

func (l locator) log(level global.Level, message string, fields global.Fields) {
	fields["commit"] = l.Commit
	fields["branch"] = l.Branch
	fields["location"] = l.RemoteURL()

	global.LogI(level, message, fields)
}
