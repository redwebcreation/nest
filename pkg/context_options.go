package pkg

import (
	"io"
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

func WithServerConfiguration(serverConfig *ServerConfiguration) ContextOption {
	return func(ctx *Context) error {
		ctx.serverConfig = serverConfig

		return nil
	}
}
