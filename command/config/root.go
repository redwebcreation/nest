package config

import (
	"github.com/redwebcreation/nest/command"
	"github.com/spf13/cobra"
)

var commands = []*cobra.Command{
	NewUseCommand(),
	NewShowCommand(),
	NewPullCommand(),
}

func NewRootConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "manage the configuration",
	}

	for _, c := range commands {
		command.Configure(c)

		cmd.AddCommand(c)
	}

	return cmd
}
