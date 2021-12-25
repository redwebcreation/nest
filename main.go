package main

import (
	"fmt"
	"os"

	"github.com/me/nest/cli"
	"github.com/me/nest/cli/proxy"
	"github.com/me/nest/global"
	"github.com/spf13/cobra"
)

func main() {
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("Nest (%s) \n", global.Version)
			os.Exit(0)
		}
	}
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

	nest.PersistentFlags().BoolP("version", "v", false, "Print version information")

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
