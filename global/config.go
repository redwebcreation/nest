package global

import (
	"github.com/mitchellh/go-homedir"
	"os"
)

// ConfigHome is the path to the global configuration for nest.
//
// It contains everything related to nest: downloaded configurations, logs, certificates, etc.
var ConfigHome string

// ConfigStoreDir is the path where the configuration repos are stored.
//
// For a given repo, the configuration is stored in the following location:
// ConfigStoreDir/<owner>-<repo>
var ConfigStoreDir string

// LogDir is the path to the log directory.
//
// It contains the following files:
// - proxy.log
// - internal.log
//
// todo: implement log rotation
var LogDir string

// CertsDir is the path to the directory containing the certificates.
// This directory should NEVER be directly accessed by the user.
// If you temper with the contents of this directory, and autocert has to generate
// a new certificate, you may get rate limited by Let's Encrypt.
var CertsDir string

// LocatorConfigFile is the path to the locator config.
var LocatorConfigFile string

// ContainerManifestFile todo: godoc
var ContainerManifestFile string

// ProxyLogFile is the path to the log file for the proxy.
//
// The path is LogDir/proxy.log
var ProxyLogFile string

// InternalLogFile is the path to the log file for the internal events.
// This is used for debugging nest itself.
// Events logged here are things such as :
// - cloned the configuration
// - containers created/removed
// - command called
// - errors, all of them.
//
// The path is LogDir/internal.log
var InternalLogFile string

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	ConfigHome = home + "/.nest"

	ConfigStoreDir = ConfigHome + "/config_store"
	CertsDir = ConfigHome + "/certs"
	LogDir = ConfigHome + "/logs"

	LocatorConfigFile = ConfigHome + "/locator.json"
	ContainerManifestFile = ConfigHome + "/manifest.json"
	ProxyLogFile = LogDir + "/proxy.log"
	InternalLogFile = LogDir + "/internal.log"

	directories := []string{ConfigStoreDir, CertsDir, LogDir}

	for _, directory := range directories {
		err = os.MkdirAll(directory, 0700)
		if err != nil {
			panic(err)
		}
	}
}
