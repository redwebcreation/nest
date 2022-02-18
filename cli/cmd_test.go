package cli

import (
	"bytes"
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

type CommandTest struct {
	Expectations   func(console *expect.Console)
	NewCommand     func(ctx *pkg.Context) (*cobra.Command, error)
	ContextOptions []pkg.ContextOption
}

func (c CommandTest) Run(t *testing.T) *pkg.Context {
	dir, err := os.MkdirTemp("", "nest-home")
	assert.NilError(t, err)

	global.configHome = dir

	// Multiplex output to a buffer as well for the raw bytes.
	buf := new(bytes.Buffer)
	console, _, err := vt10x.NewVT10XConsole(expect.WithStdout(buf))
	assert.NilError(t, err)
	defer console.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)

		c.Expectations(console)
	}()

	options := append(c.ContextOptions, pkg.WithStdio(console.Tty(), console.Tty(), console.Tty()))
	ctx, err := pkg.NewContext(options...)
	assert.NilError(t, err)

	cmd, err := c.NewCommand(ctx)
	assert.NilError(t, err)
	cmd.SetArgs([]string{})

	err = cmd.Execute()
	assert.NilError(t, err)

	// Close the slave end of the pty, and read the remaining bytes from the master end.
	console.Tty().Close()
	<-donec

	return ctx
}

type ConsoleError struct {
	err error
}

func Err(_ interface{}, err error) ConsoleError {
	return ConsoleError{err}
}

func (c ConsoleError) Check(t *testing.T) {
	// this error is expected??
	// at least it does not change the outcome of the test
	// see https://github.com/creack/pty/issues/21
	if c.err != nil && c.err.Error() == "read /dev/ptmx: input/output error" {
		return
	}

	assert.NilError(t, c.err)
}
