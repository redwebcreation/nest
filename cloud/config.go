package cloud

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/go-homedir"
	"os"
)

var ConfigDir string
var ConfigFile string

var Config *Configuration

type Configuration struct {
	Tokens []string
}

func init() {
	home, err := homedir.Dir()
	check(err)

	ConfigDir = home + "/.nestc"
	ConfigFile = ConfigDir + "/config.json"

	if _, err = os.Stat(ConfigDir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(ConfigDir, 0700)
		check(err)
	}

	if _, err = os.Stat(ConfigFile); errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(ConfigFile, []byte("{}"), 0600)
		check(err)
	}

	contents, err := os.ReadFile(ConfigFile)
	check(err)

	err = json.Unmarshal(contents, &Config)

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *Configuration) Save() error {
	contents, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(ConfigFile, contents, 0600)
	if err != nil {
		return err
	}

	return nil
}
