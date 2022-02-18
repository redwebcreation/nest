package pkg

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	logger2 "github.com/redwebcreation/nest/pkg/loggy"
	"github.com/redwebcreation/nest/pkg/manifest"
	"io"
	"log"
	"os"
	"strings"
)

type ContextOption func(*Context) error

// FileWriter provides a minimal interface for Stdin.
type FileWriter interface {
	io.Writer
	Fd() uintptr
}

// FileReader provides a minimal interface for Stdout.
type FileReader interface {
	io.Reader
	Fd() uintptr
}

func WithConfig(config *Config) ContextOption {
	return func(ctx *Context) error {
		ctx.config = config

		return nil
	}
}

func WithStdio(stdin FileReader, stdout FileWriter, stderr io.Writer) ContextOption {
	return func(ctx *Context) error {
		ctx.in = stdin
		ctx.out = stdout
		ctx.err = stderr

		return nil
	}
}

func WithServerConfig(serverConfig *ServerConfig) ContextOption {
	return func(ctx *Context) error {
		ctx.serverConfig = serverConfig

		return nil
	}
}

func WithDefaultConfigHome() ContextOption {
	return func(context *Context) error {
		for k, arg := range os.Args {
			if arg != "--config" && arg != "-c" {
				continue
			}

			if len(os.Args) <= k+1 {
				fmt.Fprintln(os.Stderr, "--config requires an argument")
				os.Exit(1)
			}

			context.home = strings.TrimRight(os.Args[k+1], "/")
			return nil
		}

		if os.Getenv("NEST_HOME") != "" {
			context.home = strings.TrimRight(os.Getenv("NEST_HOME"), "/")
			return nil
		}

		// otherwise, use the default
		userHome, err := homedir.Dir()
		if err != nil {
			return err
		}

		context.home = userHome + "/.nest"

		return nil
	}
}

func WithConfigHome(home string) ContextOption {
	return func(context *Context) error {
		context.home = home
		return nil
	}
}

func WithDefaultInternalLogger() ContextOption {
	return func(context *Context) error {
		context.logger = log.New(&logger2.FileLogger{
			Path: context.LogFile(),
		}, "", 0)

		return nil
	}
}

func WithDefaultProxyLogger() ContextOption {
	return func(context *Context) error {
		context.proxyLogger = log.New(logger2.CompositeLogger{
			Loggers: []io.Writer{
				&logger2.FileLogger{
					Path: context.ProxyLogFile(),
				},
				&logger2.FileLogger{
					Writer: os.Stdout,
				},
			},
		}, "", 0)

		return nil
	}
}

func WithManifestManager(manifestManager *manifest.Manager) ContextOption {
	return func(context *Context) error {
		context.manifestManager = manifestManager
		return nil
	}
}

func WithLogger(logger *log.Logger) ContextOption {
	return func(context *Context) error {
		context.logger = logger
		return nil
	}
}
