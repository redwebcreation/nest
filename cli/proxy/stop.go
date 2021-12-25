package proxy

import (
	"fmt"

	"github.com/spf13/cobra"
)

func stopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a proxy",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("ran: stop")
			return nil
		},
	}

	return cmd
}
