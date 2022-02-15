package main

import (
	"github.com/redwebcreation/nest/command"
	"github.com/redwebcreation/nest/command/cloud"
	"github.com/redwebcreation/nest/command/config"
	"github.com/redwebcreation/nest/command/proxy"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

var nest = &cobra.Command{
	Use:   "nest",
	Short: "Service orchestrator",
	Long:  "Nest is a powerful service orchestrator for a single server.",
}

var commands = []*cobra.Command{
	cloud.NewRootConfigCommand(),
	proxy.NewRootProxyCommand(),
	config.NewRootConfigCommand(),
	command.NewDeployCommand(),
	command.NewMedicCommand(),
	command.NewSetupCommand(),
	command.NewVersionCommand(),
	command.NewSelfUpdateCommand(),
}

func main() {
	for _, cmd := range commands {
		command.Configure(cmd)

		nest.AddCommand(cmd)
	}

	nest.PersistentFlags().StringVarP(&global.ConfigHome, "config", "c", global.ConfigHome, "set the global config path")

	err := nest.Execute()
	if err != nil {
		global.LogI(global.LevelError, err.Error(), nil)
		pkg.Stderr.Fatal("%s\n", err)
	}
}

func init() {
	nest.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})
}
