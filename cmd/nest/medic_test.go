package nest

import (
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/config"
	"github.com/redwebcreation/nest/context"
	"github.com/spf13/cobra"
	"testing"
)

func TestNewMedicCommand(t *testing.T) {
	_ = CommandTest{
		Test: func(console *expect.Console) {
			Err(console.ExpectString("Errors:")).Check(t)
			Err(console.ExpectString("- no errors")).Check(t)
			Err(console.ExpectString("Warnings:")).Check(t)
			Err(console.ExpectString("- no warnings")).Check(t)
		},
		ContextBuilder: []context.Option{
			// As the config is not nil, the context does not try to create it
			context.WithConfig(&config.Config{}),
			context.WithServicesConfig(&config.ServicesConfig{}),
		},
		NewCommand: func(ctx *context.Context) (*cobra.Command, error) {
			return NewMedicCommand(ctx), nil
		},
	}.Run(t)
}
