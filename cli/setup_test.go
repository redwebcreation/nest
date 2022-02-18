package cli

import (
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewSetupCommand(t *testing.T) {
	ctx := CommandTest{
		Test: func(console *expect.Console) {
			Err(console.ExpectString("Select your provider:")).Check(t)
			Err(console.SendLine("")).Check(t)
			Err(console.ExpectString("Enter your repository:")).Check(t)
			Err(console.SendLine("redwebcreation/nest-configs")).Check(t)
			Err(console.ExpectString("Enter your branch:")).Check(t)
			Err(console.SendLine("empty-config")).Check(t)
			Err(console.ExpectEOF()).Check(t)
		},
		NewCommand: func(ctx *pkg.Context) (*cobra.Command, error) {
			return NewSetupCommand(ctx), nil
		},
	}.Run(t)

	config, err := ctx.Config()
	assert.NilError(t, err)

	assert.Equal(t, "redwebcreation/nest-configs", config.Repository)
	assert.Equal(t, "empty-config", config.Branch)
	assert.Equal(t, "github", config.Provider)
}
