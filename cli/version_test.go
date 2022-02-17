package cli

import (
	"fmt"
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"testing"
)

func TestNewVersionCommand(t *testing.T) {
	_ = CommandTest{
		Expectations: func(console *expect.Console) {
			Err(console.ExpectString(fmt.Sprintf("Nest version %s, build %s\n", global.Version, global.Commit))).Check(t)
			Err(console.ExpectEOF()).Check(t)
		},
		NewCommand: func(ctx *pkg.Context) (*cobra.Command, error) {
			return NewVersionCommand(ctx), nil
		},
	}.Run(t)
}
