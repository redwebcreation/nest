package proxy

import (
	"fmt"

	"github.com/spf13/cobra"
)

func statusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Status of a proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("ran: status")
			return nil
		},
	}

	return cmd
}
