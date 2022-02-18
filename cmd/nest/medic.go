package nest

import (
	"fmt"
	"github.com/redwebcreation/nest/config/medic"
	"github.com/redwebcreation/nest/context"
	"github.com/spf13/cobra"
)

func runMedicCommand(ctx *context.Context) error {
	sc, err := ctx.ServerConfig()
	if err != nil {
		return err
	}

	diagnostic := medic.DiagnoseConfig(sc)

	fmt.Fprintln(ctx.Out())
	fmt.Fprintln(ctx.Out(), "Errors:")

	if len(diagnostic.Errors) == 0 {
		fmt.Fprintln(ctx.Out(), "  - no errors")
	} else {
		for _, err := range diagnostic.Errors {
			fmt.Fprintf(ctx.Out(), "  -  %s\n", err.Title)
			if err.Error != nil {
				fmt.Fprintf(ctx.Out(), "    %s\n", err.Error)
			}
		}
	}

	fmt.Fprintln(ctx.Out(), "\nWarnings:")

	if len(diagnostic.Warnings) == 0 {
		fmt.Fprintln(ctx.Out(), "  - no warnings")
	} else {
		for _, warn := range diagnostic.Warnings {
			fmt.Fprintf(ctx.Out(), "  -  %s\n", warn.Title)
			if warn.Advice != "" {
				fmt.Fprintf(ctx.Out(), "    %s\n", warn.Advice)
			}
		}
	}

	return nil
}

// NewMedicCommand creates a new `medic` command.
func NewMedicCommand(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "medic",
		Short: "diagnose your config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMedicCommand(ctx)
		},
	}

	return cmd
}
