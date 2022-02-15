package command

import (
	"fmt"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

var onlyErrors bool
var onlyWarnings bool

func runMedicCommand(cmd *cobra.Command, args []string) error {
	diagnostic := pkg.DiagnoseConfiguration()

	if !onlyWarnings {
		fmt.Println()
		pkg.Printf(pkg.Red, "  Errors:\n")

		if len(diagnostic.Errors) == 0 {
			pkg.Printf(pkg.Gray, "  - no errors\n")
		} else {
			for _, err := range diagnostic.Errors {
				pkg.Printf(pkg.White, "  -  %s\n", err.Title)
				if err.Error != nil {
					pkg.Printf(pkg.Gray, "    %s\n", err)
				}
			}
		}
	}

	if onlyErrors {
		return nil
	}

	pkg.Printf(pkg.Yellow, "\n  Warnings:")
	fmt.Println()

	if len(diagnostic.Warnings) == 0 {
		pkg.Printf(pkg.Gray, "  - no warnings\n")
	} else {
		for _, warn := range diagnostic.Warnings {
			pkg.Printf(pkg.White, "  -  %s\n", warn.Title)
			if warn.Advice != "" {
				pkg.Printf(pkg.Gray, "    %s\n", warn.Advice)
			}
		}
	}

	return nil
}

// NewMedicCommand analyses the current configuration and returns a list of errors and recommendations
func NewMedicCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "medic",
		Short: "diagnose your configuration",
		RunE:  runMedicCommand,
		Annotations: map[string]string{
			"medic": "false",
		},
	}

	cmd.Flags().BoolVarP(&onlyErrors, "only-errors", "e", false, "only show errors")
	cmd.Flags().BoolVarP(&onlyWarnings, "only-warnings", "w", false, "only show warnings")

	return cmd
}
