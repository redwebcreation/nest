package cli

import (
	"fmt"
	"github.com/redwebcreation/nest/pkg"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
)

func runMedicCommand(cmd *cobra.Command, args []string) error {
	diagnostic := pkg.DiagnoseConfiguration()

	fmt.Println()
	util.Printf(util.Red, "  Errors:\n")

	if len(diagnostic.Errors) == 0 {
		util.Printf(util.Gray, "  - no errors\n")
	} else {
		for _, err := range diagnostic.Errors {
			util.Printf(util.White, "  -  %s\n", err.Title)
			if err.Error != nil {
				util.Printf(util.Gray, "    %s\n", err)
			}
		}
	}

	util.Printf(util.Yellow, "\n  Warnings:")
	fmt.Println()

	if len(diagnostic.Warnings) == 0 {
		util.Printf(util.Gray, "  - no warnings\n")
	} else {
		for _, warn := range diagnostic.Warnings {
			util.Printf(util.White, "  -  %s\n", warn.Title)
			if warn.Advice != "" {
				util.Printf(util.Gray, "    %s\n", warn.Advice)
			}
		}
	}

	return nil
}

// NewMedicCommand creates a new `medic` command.
func NewMedicCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "medic",
		Short: "diagnose your configuration",
		RunE:  runMedicCommand,
		Annotations: map[string]string{
			"medic": "false",
		},
	}

	return cmd
}
