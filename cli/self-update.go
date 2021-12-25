package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func SelfUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self-update",
		Short: "Update the CLI to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("ran: self-update")
			return nil
		},
	}

	return cmd
}
