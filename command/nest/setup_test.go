package nest

import (
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/command"
	"github.com/redwebcreation/nest/context"
	"github.com/spf13/cobra"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewSetupCommand(t *testing.T) {
	ctx := command.CommandTest{
		Test: func(console *expect.Console) {
			command.Err(console.ExpectString("Select your provider:")).Check(t)
			command.Err(console.SendLine("")).Check(t)
			command.Err(console.ExpectString("Enter your repository:")).Check(t)
			command.Err(console.SendLine("redwebcreation/nest-configs")).Check(t)
			command.Err(console.ExpectString("Enter your branch:")).Check(t)
			command.Err(console.SendLine("empty-config")).Check(t)
			command.Err(console.ExpectEOF()).Check(t)
		},
		NewCommand: func(ctx *context.Context) (*cobra.Command, error) {
			return NewSetupCommand(ctx), nil
		},
	}.Run(t)

	config, err := ctx.Config()
	assert.NilError(t, err)

	assert.Equal(t, "redwebcreation/nest-configs", config.Repository)
	assert.Equal(t, "empty-config", config.Branch)
	assert.Equal(t, "github", config.Provider)
}
