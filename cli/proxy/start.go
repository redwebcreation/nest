package proxy

import (
	"fmt"

	"github.com/spf13/cobra"
)

func startCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start a proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("ran: start")
			return nil
		},
	}

	return cmd
}
