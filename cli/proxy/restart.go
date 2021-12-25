package proxy

import (
	"fmt"

	"github.com/spf13/cobra"
)

func restartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart a proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("ran: restart")
			return nil
		},
	}

	return cmd
}
