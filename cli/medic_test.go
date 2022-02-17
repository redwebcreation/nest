package cli

import (
	"github.com/Netflix/go-expect"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"testing"
)

func TestNewMedicCommand(t *testing.T) {
	_ = CommandTest{
		Expectations: func(console *expect.Console) {
			Err(console.ExpectString("Errors:")).Check(t)
			Err(console.ExpectString("- no errors")).Check(t)
			Err(console.ExpectString("Warnings:")).Check(t)
			Err(console.ExpectString("- no warnings")).Check(t)
		},
		ContextOptions: []pkg.ContextOption{
			// As the config is not nil, the context does not try to create it
			pkg.WithConfig(&pkg.Config{}),
			pkg.WithServerConfiguration(&pkg.ServerConfiguration{}),
		},
		NewCommand: func(ctx *pkg.Context) (*cobra.Command, error) {
			return NewMedicCommand(ctx), nil
		},
	}.Run(t)
}
