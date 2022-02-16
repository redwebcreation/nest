package config

import (
	"github.com/spf13/cobra"
)

var commands = []*cobra.Command{
	NewUseCommand(),
	NewShowCommand(),
	NewPullCommand(),
}

// NewRootCommand creates a new `config` command.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "config",
		Short: "manage the configuration",
	}

	for _, command := range commands {
		root.AddCommand(command)
	}

	return root
}
