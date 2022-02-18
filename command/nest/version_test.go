package nest

import (
	"fmt"
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/build"
	"github.com/redwebcreation/nest/command"
	"github.com/redwebcreation/nest/context"
	"github.com/spf13/cobra"
	"testing"
)

func TestNewVersionCommand(t *testing.T) {
	_ = command.CommandTest{
		Test: func(console *expect.Console) {
			command.Err(console.ExpectString(fmt.Sprintf("Nest version %s, build %s\n", build.Version, build.Commit))).Check(t)
			command.Err(console.ExpectEOF()).Check(t)
		},
		NewCommand: func(ctx *context.Context) (*cobra.Command, error) {
			return NewVersionCommand(ctx), nil
		},
	}.Run(t)
}
