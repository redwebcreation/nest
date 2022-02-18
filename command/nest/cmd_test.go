package nest

import (
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/redwebcreation/nest/context"
	"github.com/spf13/cobra"
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

type CommandTest struct {
	Test           func(console *expect.Console)
	NewCommand     func(ctx *context.Context) (*cobra.Command, error)
	ContextBuilder []context.ContextOption
	Setup          func(ctx *context.Context) []context.ContextOption
}

func (c CommandTest) Run(t *testing.T) *context.Context {
	dir, err := os.MkdirTemp("", "nest-home")
	assert.NilError(t, err)

	console, _, err := vt10x.NewVT10XConsole()
	assert.NilError(t, err)

	defer console.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)

		c.Test(console)
	}()

	ctx, err := context.NewContext(context.WithConfigHome(dir), context.WithStdio(console.Tty(), console.Tty(), console.Tty()))
	assert.NilError(t, err)

	for _, option := range c.ContextBuilder {
		err = option(ctx)
		assert.NilError(t, err)
	}

	if c.Setup != nil {
		for _, opt := range c.Setup(ctx) {
			err = opt(ctx)
			assert.NilError(t, err)
		}
	}

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
