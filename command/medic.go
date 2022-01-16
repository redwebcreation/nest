package command

import (
	"encoding/json"
	"fmt"
	"github.com/redwebcreation/nest/util"

	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

var jsonFormat bool
var onlyErrors bool
var onlyWarnings bool

func runMedicCommand(cmd *cobra.Command, args []string) error {
	diagnosis := pkg.DiagnoseConfiguration()

	if jsonFormat {
		out, _ := json.Marshal(diagnosis)

		fmt.Println(string(out))
		return nil
	}

	if !onlyWarnings {
		fmt.Println()
		fmt.Printf("  %sErrors:%s\n", util.Red, util.Reset())

		if len(diagnosis.Errors) == 0 {
			fmt.Printf("  %s- no errors%s\n", util.Gray, util.Reset())
		} else {
			for _, err := range diagnosis.Errors {
				fmt.Printf("  %s- %s%s\n", util.White, err.Title, util.Reset())
				if err.Error != nil {
					fmt.Printf("    %s%s%s\n", util.Gray, err.Error, util.Reset())
				}
			}
		}
	}

	if !onlyErrors {
		fmt.Printf("\n  %sWarnings:%s\n", util.Yellow, util.Reset())

		if len(diagnosis.Warnings) == 0 {
			fmt.Printf("  %s- no warnings%s\n", util.Gray, util.Reset())
		} else {
			for _, warn := range diagnosis.Warnings {
				fmt.Printf("  %s- %s%s\n", util.White, warn.Title, util.Reset())
				if warn.Advice != "" {
					fmt.Printf("    %s%s%s\n", util.Gray, warn.Advice, util.Reset())
				}
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

	cmd.Flags().BoolVarP(&jsonFormat, "jsonFormat", "j", false, "output in jsonFormat format")
	cmd.Flags().BoolVarP(&onlyErrors, "only-errors", "e", false, "only show errors")
	cmd.Flags().BoolVarP(&onlyWarnings, "only-warnings", "w", false, "only show warnings")

	return cmd
}
