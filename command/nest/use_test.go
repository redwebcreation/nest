package nest

import (
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/command"
	"github.com/redwebcreation/nest/config"
	"github.com/redwebcreation/nest/context"
	"github.com/spf13/cobra"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewUseCommand(t *testing.T) {
	ctx := command.CommandTest{
		Test: func(console *expect.Console) {
			command.Err(console.SendLine("")).Check(t)
			command.Err(console.ExpectEOF()).Check(t)
		},
		Setup: func(ctx *context.Context) []context.ContextOption {
			return []context.ContextOption{
				context.WithConfig(&config.Config{
					Provider:   "github",
					Repository: "redwebcreation/nest-configs",
					Branch:     "empty-config",
					Path:       ctx.ConfigFile(),
					StoreDir:   ctx.ConfigStoreDir(),
					Logger:     ctx.Logger(),
					Git: &config.GitWrapper{
						Logger: ctx.Logger(),
					},
				}),
			}
		},
		NewCommand: func(ctx *context.Context) (*cobra.Command, error) {
			config, err := ctx.Config()
			if err != nil {
				return nil, err
			}

			err = config.Clone()
			if err != nil {
				return nil, err
			}

			return NewUseCommand(ctx), nil
		},
	}.Run(t)

	config, err := ctx.Config()
	assert.NilError(t, err)

	// see https://github.com/redwebcreation/nest-configs/tree/empty-config
	assert.Equal(t, config.Commit, "3ea941eaf6d2bfcc97480ce5df49bee91d8f09e2")
}
