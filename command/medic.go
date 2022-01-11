package command

import (
	"fmt"

	"github.com/redwebcreation/nest/common"
	"github.com/spf13/cobra"
)

func runMedicCommand(cmd *cobra.Command, args []string) error {
	diagnosis := common.DiagnoseConfiguration()

	fmt.Printf("\n  \033[1m\033[38;2;255;60;0mErrors:\033[0m\n")

	if len(diagnosis.Errors) == 0 {
		fmt.Println("  \033[1m\033[38;2;120;120;120m- no errors\033[0m")
	} else {
		for _, err := range diagnosis.Errors {
			fmt.Printf("  \033[1m\033[38;2;255;255;255m- %s\033[0m\n", err.Title)
			if err.Error != nil {
				fmt.Printf("  \033[1m\033[38;2;125;125;125m  %s\033[0m\n", err.Error.Error())
			}
		}
	}

	fmt.Printf("\n  \033[1m\033[38;2;250;175;0mRecommendations:\033[0m\n")

	if len(diagnosis.Warnings) == 0 {
		fmt.Println("  \033[1m\033[38;2;120;120;120m- no recommendations\033[0m")
	} else {
		for _, warn := range diagnosis.Warnings {
			fmt.Printf("  \033[1m\033[38;2;255;255;255m- %s\033[0m\n", warn.Title)
			if warn.Advice != "" {
				fmt.Printf("  \033[1m\033[38;2;125;125;125m  %s\033[0m\n", warn.Advice)
			}
		}
	}

	return nil
}

func MedicCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "medic",
		Short: "diagnose your configuration",
		RunE:  runMedicCommand,
	}

	return WithConfig(cmd)
}
