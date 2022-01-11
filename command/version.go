package command

import (
	"fmt"

	"github.com/redwebcreation/nest/global"
	"github.com/spf13/cobra"
)

func runVersionCommand(cmd *cobra.Command, args []string) error {
	fmt.Printf("nest@%s\n", global.Version)
	return nil
}

// NewVersionCommand prints nest's version
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print nest's version",
		RunE:  runVersionCommand,
	}

	return cmd
}
