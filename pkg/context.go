package pkg

import (
	"github.com/redwebcreation/nest/global"
	"io"
	"os"
)

type Context struct {
	config       *Config
	serverConfig *ServerConfiguration
	out          FileWriter
	in           FileReader
	err          io.Writer
}

func (c *Context) HasConfig() bool {
	if c.serverConfig != nil {
		return true
	}

	_, err := os.ReadFile(global.ConfigFile())

	return err == nil
}

func (c *Context) Config() (*Config, error) {
	if c.config == nil {
		config, err := NewConfig()
		if err != nil {
			return nil, err
		}

		c.config = config
	}

	return c.config, nil
}

func (c *Context) ServerConfiguration() (*ServerConfiguration, error) {
	config, err := c.Config()
	if err != nil {
		return nil, err
	}

	if c.serverConfig == nil {
		serverConfig, err := config.GetServerConfiguration()
		if err != nil {
			return nil, err
		}

		c.serverConfig = serverConfig
	}

	err = DiagnoseConfiguration(c.serverConfig).MustPass()
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

func NewContext(opts ...ContextOption) (*Context, error) {
	ctx := &Context{}
	for _, opt := range opts {
		if err := opt(ctx); err != nil {
			return nil, err
		}
	}
	return ctx, nil
}
