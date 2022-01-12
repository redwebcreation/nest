package command

import (
	"fmt"
	"github.com/redwebcreation/nest/util"

	"github.com/redwebcreation/nest/common"
	"github.com/spf13/cobra"
)

func runMedicCommand(cmd *cobra.Command, args []string) error {
	diagnosis := common.DiagnoseConfiguration()

	fmt.Println()
	fmt.Printf("  %sErrors:%s\n", util.Red, util.Reset)

	if len(diagnosis.Errors) == 0 {
		fmt.Printf("  %s- no errors%s", util.Gray, util.Reset)
	} else {
		for _, err := range diagnosis.Errors {
			fmt.Printf("  %s- %s%s\n", util.White, err, util.Reset)
			if err.Error != nil {
				fmt.Printf("    %s%s%s\n", util.Gray, err.Error, util.Reset)
			}
		}
	}

	fmt.Printf("\n\n  %sWarnings:%s\n", util.Yellow, util.Reset)

	if len(diagnosis.Warnings) == 0 {
		fmt.Printf("  %s- no warnings%s", util.Gray, util.Reset)
	} else {
		for _, warn := range diagnosis.Warnings {
			fmt.Printf("  %s- %s%s\n", util.White, warn.Title, util.Reset)
			if warn.Advice != "" {
				fmt.Printf("    %s%s%s\n", util.Gray, warn.Advice, util.Reset)
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
	}

	return cmd
}
