package main

import (
	"fmt"
	"github.com/me/nest/cli"
	"github.com/me/nest/cli/proxy"
	"github.com/me/nest/common"
	"github.com/me/nest/global"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("nest@%s\n", global.Version)
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
		cli.MedicCommand(),
		cli.SelfUpdateCommand(),
		cli.ReconfigureCommand(),
	}

	for _, command := range commands {
		command.SilenceUsage = true
		command.SilenceErrors = true

		nest.AddCommand(command)
	}

	nest.PersistentFlags().BoolP("version", "v", false, "print version information")

	// hide the help command
	nest.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})

	nest.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "medic" {
			return nil
		}

		diagnosis := common.DiagnoseConfiguration()

		if len(diagnosis.Errors) == 0 {
			return nil
		}

		return fmt.Errorf("your configuration is invalid, please run `nest medic` for details")
	}

	nest.SilenceErrors = true

	err := nest.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error: "+err.Error())
		os.Exit(1)
	}
}
