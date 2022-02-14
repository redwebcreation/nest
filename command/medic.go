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
		pkg.Red.Render("  Errors:")

		if len(diagnostic.Errors) == 0 {
			pkg.Gray.Render("  - no errors")
		} else {
			for _, err := range diagnostic.Errors {
				pkg.White.Render("  -  " + err.Title)
				if err.Error != nil {
					pkg.Gray.Render("    " + err.Error.Error())
				}
			}
		}
	}

	if onlyErrors {
		return nil
	}

	pkg.Yellow.Render("	Warnings:")

	if len(diagnostic.Warnings) == 0 {
		pkg.Gray.Render("  - no warnings")
	} else {
		for _, warn := range diagnostic.Warnings {
			pkg.White.Render("  -  " + warn.Title)
			if warn.Advice != "" {
				pkg.Gray.Render("    " + warn.Advice)
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
