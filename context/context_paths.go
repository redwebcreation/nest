package context

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (c *Context) ConfigStoreDir() string {
	return ensureDirExists(c.Home() + "/config-store")
}

func (c *Context) configFile() string {
	return ensureDirExists(c.Home() + "/config.json")
}

func (c *Context) manifestsDir() string {
	return ensureDirExists(c.Home() + "/manifests")
}

func (c *Context) certsDir() string {
	return ensureDirExists(c.Home() + "/certs")
}

func (c *Context) logFile() string {
	return ensureDirExists(c.Home() + "/logs/internal.log")
}

func (c *Context) proxyLogFile() string {
	return ensureDirExists(c.Home() + "/logs/proxy.log")
}

func (c *Context) cloudCredentialsFile() string {
	return c.Home() + "/.creds"
}

func (c *Context) subnetRegistryPath() string {
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
