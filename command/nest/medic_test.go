package nest

import (
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/command"
	"github.com/redwebcreation/nest/config"
	"github.com/redwebcreation/nest/context"
	"github.com/spf13/cobra"
	"testing"
)

func TestNewMedicCommand(t *testing.T) {
	_ = command.CommandTest{
		Test: func(console *expect.Console) {
			command.Err(console.ExpectString("Errors:")).Check(t)
			command.Err(console.ExpectString("- no errors")).Check(t)
			command.Err(console.ExpectString("Warnings:")).Check(t)
			command.Err(console.ExpectString("- no warnings")).Check(t)
		},
		ContextBuilder: []context.ContextOption{
			// As the config is not nil, the context does not try to create it
			context.WithConfig(&config.Config{}),
			context.WithServerConfig(&config.ServerConfig{}),
		},
		NewCommand: func(ctx *context.Context) (*cobra.Command, error) {
			return NewMedicCommand(ctx), nil
		},
	}.Run(t)
}
