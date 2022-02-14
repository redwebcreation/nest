package global

import (
	"github.com/mitchellh/go-homedir"
	"os"
)

// ConfigHome is the path to the global configuration for nest.
//
// It can be set using the NEST_HOME environment variable.
var ConfigHome string

func GetConfigStoreDir() string {
	return ConfigHome + "/configStore"
}

func GetCertsDir() string {
	return ensureExists(ConfigHome + "/certs")
}

func GetLogsDir() string {
	return ensureExists(ConfigHome + "/logs")
}

func GetLocatorConfigFile() string {
	return ConfigHome + "/locator.json"
}

func GetContainerManifestFile() string {
	return ConfigHome + "/manifest.json"
}

func GetProxyLogFile() string {
	return GetLogsDir() + "/proxy.log"
}

func GetInternalLogFile() string {
	return GetLogsDir() + "/internal.log"
}

func init() {
	if os.Getenv("NEST_HOME") != "" {
		ConfigHome = os.Getenv("NEST_HOME")
		return
	}

	// otherwise, use the default
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	ConfigHome = ensureExists(home + "/.nest")
}

func ensureExists(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	return path
}
