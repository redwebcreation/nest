package context

import (
	"github.com/c-robinson/iplib"
	"github.com/redwebcreation/nest/cloud"
	"github.com/redwebcreation/nest/config"
	"github.com/redwebcreation/nest/config/medic"
	"github.com/redwebcreation/nest/deploy"
	"github.com/redwebcreation/nest/docker"
	"golang.org/x/crypto/acme/autocert"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
)

// Context is a struct that holds the context of the application
type Context struct {
	// home is the path to the logger config for nest.
	//
	// It resolves to the following (in order):
	// - WithConfigHome option
	// - the --config/-c flag
	// - $NEST_HOME
	// - ~/.nest
	home string
	// config contains the path to nest's config file.
	// it is resolved once and cached.
	config *config.Config
	// servicesConfig contains the resolved server config from the config.
	// it is resolved once and cached.
	servicesConfig *config.ServicesConfig
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

	manifestManager *deploy.Manager
	subnetter       *docker.Subnetter
}

// Config returns the cached nest config or loads it if it hasn't been loaded yet.
func (c *Context) Config() (*config.Config, error) {
	if c.config == nil {
		cf, err := config.NewConfig(c.configFile(), c.ConfigStoreDir(), c.Logger())
		if err != nil {
			return nil, err
		}

		c.config = cf
	}

	return c.config, nil
}

func (c *Context) UnvalidatedServicesConfig() (*config.ServicesConfig, error) {
	nc, err := c.Config()
	if err != nil {
		return nil, err
	}

	if c.servicesConfig == nil {
		servicesConfig, err := nc.ServicesConfig()
		if err != nil {
			return nil, err
		}

		c.servicesConfig = servicesConfig
	}

	err = c.servicesConfig.ExpandIncludes(nc)
	if err != nil {
		return nil, err
	}

	return c.servicesConfig, nil
}

// ServicesConfig returns the cached services config or loads it if it hasn't been loaded yet.
func (c *Context) ServicesConfig() (*config.ServicesConfig, error) {
	_, err := c.UnvalidatedServicesConfig()
	if err != nil {
		return nil, err
	}

	err = medic.DiagnoseConfig(c.servicesConfig).MustPass()
	if err != nil {
		return nil, err
	}

	return c.servicesConfig, nil
}

func (c *Context) Out() FileWriter {
	if c.out == nil {
		c.out = os.Stdout
	}

	return c.out
}

func (c *Context) In() FileReader {
	if c.in == nil {
		c.in = os.Stdin
	}

	return c.in
}

func (c *Context) Err() io.Writer {
	if c.err == nil {
		c.err = os.Stderr
	}

	return c.err
}

func (c *Context) Home() string {
	return c.home
}

func (c *Context) ProxyLogger() *log.Logger {
	return c.proxyLogger
}

func (c *Context) Logger() *log.Logger {
	return c.logger
}

func (c *Context) ManifestManager() *deploy.Manager {
	if c.manifestManager == nil {
		c.manifestManager = &deploy.Manager{
			Path: c.manifestsDir(),
		}
	}

	return c.manifestManager
}

func New(opts ...Option) (*Context, error) {
	ctx := &Context{}
	defaultOptions := []Option{
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

func (c *Context) CloudCredentials() (id, token string, err error) {
	bytes, err := os.ReadFile(c.cloudCredentialsFile())
	if err != nil {
		return "", "", err
	}

	credentials := string(bytes)

	if err = cloud.ValidateDsn(credentials); err != nil {
		return "", "", err
	}

	return cloud.ParseDsn(credentials)
}

func (c *Context) CloudClient() (*cloud.Client, error) {
	id, token, err := c.CloudCredentials()
	if err != nil {
		return nil, err
	}

	return cloud.NewClient(id, token), nil
}

func (c *Context) Subnetter(CIDRs []string) *docker.Subnetter {
	if c.subnetter == nil {
		var subnets []iplib.Net4
		for _, cidr := range CIDRs {
			// todo(medic):  validate cidr
			_, n, err := iplib.ParseCIDR(cidr)
			if err != nil {
				panic(err)
			}

			_, mask := n.Mask().Size()
			subnets = append(subnets, iplib.NewNet4(n.IP(), mask))
		}

		if len(CIDRs) == 0 {
			subnets = []iplib.Net4{
				iplib.NewNet4(net.IPv4(10, 0, 0, 0), 8),
			}
		}

		c.subnetter = &docker.Subnetter{
			Lock:         &sync.Mutex{},
			RegistryPath: c.subnetRegistryPath(),
			Subnets:      subnets,
		}
	}

	return c.subnetter
}

func (c *Context) SetCloudCredentials(id string, token string) error {
	return ioutil.WriteFile(c.cloudCredentialsFile(), []byte(cloud.FormatDsn(id, token)), 0600)
}

func (c *Context) CertificateStore() autocert.DirCache {
	return autocert.DirCache(c.certsDir())
}

func (c *Context) NewConfig(provider, repository, branch string) *config.Config {
	return &config.Config{
		Provider:   provider,
		Repository: repository,
		Branch:     branch,
		Path:       c.configFile(),
		StoreDir:   c.ConfigStoreDir(),
		Logger:     c.Logger(),
		Git: &config.Git{
			Logger: c.Logger(),
		},
	}
}
