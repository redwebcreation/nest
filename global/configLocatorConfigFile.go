package global

import "github.com/mitchellh/go-homedir"

// ConfigLocatorConfigFile is the path to the config locator file
var ConfigLocatorConfigFile string

func FindConfigLocatorConfigFile() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return home + "/.nest.json", nil
}

func init() {
	configLocatorConfigFile, err := FindConfigLocatorConfigFile()
	if err != nil {
		panic(err)
	}
	ConfigLocatorConfigFile = configLocatorConfigFile
}
