package pkg

import (
	"github.com/redwebcreation/nest/pkg/manifest"
	"io"
	"log"
	"os"
)

// Context is a struct that holds the context of the application
type Context struct {
	// home is the path to the global config for nest.
	//
	// It resolves to the following (in order):
	// - WithConfigHome option
	// - the --config/-c flag
	// - $NEST_HOME
	// - ~/.nest
	home string
	// config contains the path to nest's config file.
	// it is resolved once and cached.
	config *Config
	// serverConfig contains the resolved server config from the config.
	// it is resolved once and cached.
	serverConfig *ServerConfig
	// out is a minimal interface to write to stdout.
	out FileWriter
	// in is a minimal interface to read from stdin.
	in FileReader
	// err is a minimal interface to write to stderr.
	err io.Writer
	// logger is nest's internal logger.
	// it is used to log any action that changes any kind of state.
	logger *log.Logger
	// proxyLogger is solely used to log proxy events such as a request coming in, an error in the proxy, etc.
	proxyLogger *log.Logger

	manifestManager *manifest.Manager
}

// HasConfig returns true if the serverConfig has been loaded or the config file exists.
func (c *Context) HasConfig() bool {
	if c.serverConfig != nil {
		return true
	}

	_, err := os.ReadFile(c.ConfigFile())

	return err == nil
}

// Config returns the cached nest config or loads it if it hasn't been loaded yet.
func (c *Context) Config() (*Config, error) {
	if c.config == nil {
		config, err := NewConfig(c.ConfigFile(), c.ConfigStoreDir(), c.Logger())
		if err != nil {
			return nil, err
		}

		c.config = config
	}

	return c.config, nil
}

// ServerConfig returns the cached server config or loads it if it hasn't been loaded yet.
func (c *Context) ServerConfig() (*ServerConfig, error) {
	config, err := c.Config()
	if err != nil {
		return nil, err
	}

	if c.serverConfig == nil {
		serverConfig, err := config.ServerConfig()
		if err != nil {
			return nil, err
		}

		c.serverConfig = serverConfig
	}

	err = DiagnoseConfig(c.serverConfig).MustPass()
	if err != nil {
		return nil, err
	}

	return c.serverConfig, nil
}

func (c *Context) Out() FileWriter {
	if c.out == nil {
		c.out = os.Stdout
	}

	return c.out
}

func (c Context) In() FileReader {
	if c.in == nil {
		c.in = os.Stdin
	}

	return c.in
}

func (c Context) Err() io.Writer {
	if c.err == nil {
		c.err = os.Stderr
	}

	return c.err
}

func (c Context) Home() string {
	return c.home
}

func (c Context) ProxyLogger() *log.Logger {
	return c.proxyLogger
}
func (c Context) Logger() *log.Logger {
	return c.logger
}

func (c Context) ManifestManager() *manifest.Manager {
	if c.manifestManager == nil {
		c.manifestManager = &manifest.Manager{
			Path: c.ManifestsDir(),
		}
	}

	return c.manifestManager
}

func NewContext(opts ...ContextOption) (*Context, error) {
	ctx := &Context{}
	defaultOptions := []ContextOption{
		WithDefaultConfigHome(),
		WithDefaultInternalLogger(),
		WithDefaultProxyLogger(),
	}

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(ctx); err != nil {
			return nil, err
		}
	}

	return ctx, nil
}
