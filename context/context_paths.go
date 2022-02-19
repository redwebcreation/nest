package context

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (c *Context) ManifestsDir() string {
	return ensureDirExists(c.Home() + "/manifests")
}

func (c Context) ConfigStoreDir() string {
	return ensureDirExists(c.Home() + "/config-store")
}

func (c *Context) CertsDir() string {
	return ensureDirExists(c.Home() + "/certs")
}

func (c Context) LogFile() string {
	return ensureDirExists(c.Home() + "/logs/internal.log")
}

func (c Context) ProxyLogFile() string {
	return ensureDirExists(c.Home() + "/logs/proxy.log")
}

func (c Context) ConfigFile() string {
	return ensureDirExists(c.Home() + "/config.json")
}

func (c *Context) SelfSignedKeyFile() string {
	return ensureDirExists(c.CertsDir() + "/testing_key.pem")
}

func (c *Context) SelfSignedCertFile() string {
	return ensureDirExists(c.CertsDir() + "/testing_cert.pem")
}

func (c *Context) CloudCredentialsFile() string {
	return c.Home() + "/.creds"
}

func (c *Context) SubnetRegistryPath() string {
	return ensureDirExists(c.Home() + "/subnets.list")
}

// ensureDirExists creates all the directories in a given path if they don't exist.
func ensureDirExists(path string) string {
	// if the path contains a filename, create all its parent directories
	filename := filepath.Base(path)
	var isFilename bool
	if !strings.HasPrefix(filename, ".") && strings.Contains(filename, ".") {
		path = filepath.Dir(path)
		isFilename = true
	}

	_, err := os.Stat(path)

	switch {
	case errors.Is(err, fs.ErrNotExist):
		_ = os.MkdirAll(path, 0755)
	case err != nil:
		panic(err)
	}

	if isFilename {
		path += "/" + filename
	}

	return strings.TrimRight(path, "/") // remove trailing slash
}
