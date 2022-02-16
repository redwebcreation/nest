package global

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
)

// ConfigHome is the path to the global configuration for nest.
//
// It resolves to the following (in order):
// - the --config/-c flag
// - $NEST_HOME
// - ~/.nest
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

func GetManifestsDir() string {
	return ensureExists(ConfigHome + "/manifests")
}

func GetLocatorConfigFile() string {
	return ConfigHome + "/locator.json"
}

func GetContainerManifestFile(manifest string) string {
	return GetManifestsDir() + "/" + manifest + ".json"
}

func GetProxyLogFile() string {
	return GetLogsDir() + "/proxy.log"
}

func GetInternalLogFile() string {
	return GetLogsDir() + "/internal.log"
}

func GetCloudTokenFile() string {
	return ConfigHome + "/cloudToken.json"
}

func GetSelfSignedCertKeyFile() string {
	return ConfigHome + "/certs/dev_key.pem"
}

func GetSelfSignedCertFile() string {
	return ConfigHome + "/certs/dev_cert.pem"
}

func init() {
	for k, arg := range os.Args {
		if arg != "--config" && arg != "-c" {
			continue
		}

		if len(os.Args) <= k+1 {
			fmt.Fprintln(os.Stderr, "--config requires an argument")
			os.Exit(1)
		}

		ConfigHome = os.Args[k+1]
		return
	}

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
		err = os.MkdirAll(path, 0700)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	return path
}
