package config

import (
	"fmt"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

func runShowCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("location:", pkg.Locator.RemoteURL())
	fmt.Printf("current commit: %s\n", pkg.Locator.Commit[:7])
	fmt.Println("branch:", pkg.Locator.Branch)

	// todo: list files?

	return nil
}

// NewShowCommand creates a new `show` command.
func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "prints the current configuration",
		RunE:  runShowCommand,
	}

	return cmd
}
