package global

import "github.com/mitchellh/go-homedir"

// ConfigLocatorConfigFile is the path to the config locator file
var ConfigLocatorConfigFile string

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	ConfigLocatorConfigFile = home + "/.nest.json"
}
