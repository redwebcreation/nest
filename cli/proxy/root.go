package proxy

import "github.com/spf13/cobra"

func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Proxy related commands",
	}

	cmd.AddCommand(startCommand())
	cmd.AddCommand(stopCommand())
	cmd.AddCommand(restartCommand())
	cmd.AddCommand(statusCommand())

	return cmd
}
