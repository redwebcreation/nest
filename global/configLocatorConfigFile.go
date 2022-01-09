package global

import "github.com/mitchellh/go-homedir"

var ConfigLocatorConfigFile string

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	ConfigLocatorConfigFile = home + "/.nest.json"
}
