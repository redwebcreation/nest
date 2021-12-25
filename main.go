package main

import (
	"fmt"
	"os"

	"github.com/me/nest/cli"
	"github.com/me/nest/cli/proxy"
	"github.com/spf13/cobra"
)

func main() {
	nest := &cobra.Command{
		Use:   "nest",
		Short: "Service orchestrator",
		Long:  "Nest is a powerful service orchestrator for a single server.",
	}

	commands := []*cobra.Command{
		proxy.RootCommand(),
		cli.DeployCommand(),
		cli.DiagnoseCommand(),
		cli.SelfUpdateCommand(),
	}

	for _, command := range commands {
		command.SilenceUsage = true
		command.SilenceErrors = true
		nest.AddCommand(command)
	}

	nest.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})

	err := nest.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: "+err.Error())
		os.Exit(1)
	}
}
