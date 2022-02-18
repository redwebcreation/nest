package cli

import (
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"gotest.tools/v3/assert"
	"testing"
)

func TestNewUseCommand(t *testing.T) {
	ctx := CommandTest{
		Test: func(console *expect.Console) {
			Err(console.SendLine("")).Check(t)
			Err(console.ExpectEOF()).Check(t)
		},
		Setup: func(ctx *pkg.Context) []pkg.ContextOption {
			return []pkg.ContextOption{
				pkg.WithConfig(&pkg.Config{
					Provider:   "github",
					Repository: "redwebcreation/nest-configs",
					Branch:     "empty-config",
					Path:       ctx.ConfigFile(),
					StoreDir:   ctx.ConfigStoreDir(),
					Logger:     ctx.Logger(),
					Git: &pkg.GitWrapper{
						Logger: ctx.Logger(),
					},
				}),
			}
		},
		NewCommand: func(ctx *pkg.Context) (*cobra.Command, error) {
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
