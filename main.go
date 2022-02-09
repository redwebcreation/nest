package main

import (
	"fmt"
	"github.com/redwebcreation/nest/command"
	"github.com/redwebcreation/nest/command/config"
	"github.com/redwebcreation/nest/command/proxy"
	"github.com/redwebcreation/nest/global"
	"os"

	"github.com/spf13/cobra"
)

var nest = &cobra.Command{
	Use:   "nest",
	Short: "Service orchestrator",
	Long:  "Nest is a powerful service orchestrator for a single server.",
}

var commands = []*cobra.Command{
	command.NewDeployCommand(),
	command.NewMedicCommand(),
	config.NewRootConfigCommand(),
	proxy.NewRootProxyCommand(),
	command.NewSetupCommand(),
	command.NewVersionCommand(),
	command.NewSelfUpdateCommand(),
}

func main() {
	for _, cmd := range commands {
		command.Configure(cmd)

		nest.AddCommand(cmd)
	}

	nest.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})
	nest.PersistentFlags().StringVarP(&global.ConfigHome, "config", "c", global.ConfigHome, "set the global config path")

	err := nest.Execute()
	if err != nil {
		global.InternalLogger.Log(
			global.LevelError,
			"main",
			global.NewField("error", err.Error()),
		)
		_, _ = fmt.Fprintln(os.Stderr, "error: "+err.Error())
		os.Exit(1)
	}
}
