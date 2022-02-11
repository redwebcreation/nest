package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/util"
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
	Secrets    map[string][]byte
	VCS        *util.VCS `yaml:"-"`
}

func (l locator) ConfigPath() string {
	return global.ConfigStoreDir + "/" + strings.Replace(l.Repository, "/", "-", -1)
}
func (l locator) RemoteURL() string {
	return fmt.Sprintf("git@%s.com:%s.git", l.Provider, l.Repository)
}

func (l locator) Read(file string) ([]byte, error) {
	configPath := l.ConfigPath()

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		err = l.VCS.Clone(l.RemoteURL(), l.ConfigPath())

		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return l.VCS.ReadFile(configPath, l.Commit, file)
}

func (l locator) Resolve() (*Configuration, error) {
	err := l.Load()
	if err != nil {
		return nil, err
	}

	contents, err := l.Read("nest.yml")
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

func (l locator) Validate() error {
	return nil
}

func (l *locator) Load() error {
	contents, err := os.ReadFile(global.LocatorConfigFile)
	if err != nil {
		return err
	}

	var p locator
	err = json.Unmarshal(contents, &p)
	if err != nil {
		return err
	}

	l.Provider = p.Provider
	l.Repository = p.Repository
	l.Branch = p.Branch
	l.Commit = p.Commit

	if l.Commit == "" {
		return fmt.Errorf("commit is empty, run `nest setup` to set it")
	}

	return nil
}

func (l *locator) LoadCommit(commit string) error {
	l.Commit = commit
	return l.Load()
}

func init() {
	Locator.VCS = util.VcsGit
}
