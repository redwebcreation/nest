package main

import (
	"fmt"
	"os"

	"github.com/me/nest/cli"
	"github.com/me/nest/cli/proxy"
	"github.com/me/nest/common"
	"github.com/me/nest/global"
	"github.com/spf13/cobra"
)

func main() {

	//os.Exit(0)
	nest := &cobra.Command{
		Use:   "nest",
		Short: "Service orchestrator",
		Long:  "Nest is a powerful service orchestrator for a single server.",
	}

	for _, command := range []*cobra.Command{
		proxy.RootCommand(),
		cli.DeployCommand(),
		cli.MedicCommand(),
		cli.ConfigCommand(),
		cli.SelfUpdateCommand(),
		cli.ConfigureCommand(),
		cli.VersionCommand(),
	} {
		command.SilenceUsage = true
		command.SilenceErrors = true

		nest.AddCommand(command)
	}

	// hide the help command
	nest.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})

	nest.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		commandName := cmd.Name()
		if !global.IsConfigLocatorConfigured && (commandName != "configure" && commandName != "version" && commandName != "self-update") {
			return fmt.Errorf("run `nest configure` to set up nest")
		}

		if commandName == "medic" || commandName == "version" || commandName == "configure" || commandName == "self-update" {
			return nil
		}

		diagnosis := common.DiagnoseConfiguration()

		if len(diagnosis.Errors) == 0 {
			return nil
		}

		return fmt.Errorf("your configuration is invalid, please run `nest medic` to troubleshoot")
	}

	nest.SilenceErrors = true

	err := nest.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error: "+err.Error())
		os.Exit(1)
	}
}
