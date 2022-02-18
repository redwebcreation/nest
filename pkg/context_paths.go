package pkg

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func (c *Context) ManifestsDir() string {
	return ensureExists(c.Home() + "/manifests")
}

func (c Context) ConfigStoreDir() string {
	return ensureExists(c.Home() + "/config-store")
}

func (c *Context) CertsDir() string {
	return ensureExists(c.Home() + "/certs")
}

func (c Context) LogFile() string {
	return ensureExists(c.Home() + "/logs/internal.log")
}

func (c Context) ProxyLogFile() string {
	return ensureExists(c.Home() + "/logs/proxy.log")
}

func (c Context) ConfigFile() string {
	return ensureExists(c.Home() + "/config.json")
}

func (c *Context) SelfSignedKeyFile() string {
	return ensureExists(c.CertsDir() + "/testing_key.pem")
}

func (c *Context) SelfSignedCertFile() string {
	return ensureExists(c.CertsDir() + "/testing_cert.pem")
}

func (c *Context) CloudTokenFile() string {
	return ensureExists(c.Home() + "/cloud-token")
}

func (c *Context) ManifestFile(id string) string {
	return ensureExists(c.ManifestsDir() + "/" + id + ".json")
}

// ensureExists creates all the directories in a given path if they don't exist.
func ensureExists(path string) string {
	// if the path contains a filename, create all its parent directories
	filename := filepath.Base(path)
	if !strings.HasPrefix(filename, ".") && strings.Contains(filename, ".") {
		path = filepath.Dir(path)
	}

	_, err := os.Stat(path)

	switch {
	case errors.Is(err, fs.ErrNotExist):
		_ = os.MkdirAll(path, 0755)
	case err != nil:
		panic(err)
	}

	return strings.TrimRight(path, "/") // remove trailing slash
}
