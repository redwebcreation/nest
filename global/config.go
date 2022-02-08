package global

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"os"
)

// ConfigHome is the path to the global configuration for nest.
var ConfigHome string

// LocatorConfigFile is the path to the locator config.
var LocatorConfigFile string

var ConfigStoreDir string

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	ConfigHome = home + "/.nest"
	LocatorConfigFile = ConfigHome + "/locator.json"
	ConfigStoreDir = ConfigHome + "/git_store"

	// Create the config directory if it doesn't exist.
	if _, err = os.Stat(ConfigHome); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(ConfigHome, 0700)
		if err != nil {
			panic(err)
		}
	}

	directories := []string{
		"certs",
		"git_store",
	}

	for _, directory := range directories {
		err = os.Mkdir(ConfigHome+"/"+directory, 0700)
		if err != nil && !errors.Is(err, os.ErrExist) {
			panic(err)
		}
	}

}
