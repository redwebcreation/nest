package main

import (
	"fmt"
	"github.com/redwebcreation/nest/cli"
	"github.com/redwebcreation/nest/cli/cloud"
	"github.com/redwebcreation/nest/cli/proxy"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"os"
)

func newNestCommand(ctx *pkg.Context) *cobra.Command {
	nest := &cobra.Command{
		Use:           "nest",
		Short:         "Service orchestrator",
		SilenceUsage:  true,
		SilenceErrors: true,
		Long:          "Nest is a powerful service orchestrator for a single server.",
		Version:       fmt.Sprintf("%s, build %s", global.Version, global.Commit),
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			global.LogI(
				global.LevelDebug,
				"command invoked",
				global.Fields{
					"tag":     "command.invoke",
					"command": cmd.Name(),
				},
			)
		},
	}

	nest.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})

	// This flag is not actually used by any of the commands.
	// Its value is used in the init function in global/server.go
	nest.PersistentFlags().StringP("config", "c", global.ConfigHome, "set the global config path")

	nest.AddCommand(
		// version
		cli.NewVersionCommand(ctx),

		// setup
		cli.NewSetupCommand(ctx),

		// use
		cli.NewUseCommand(ctx),

		// medic
		cli.NewMedicCommand(ctx),

		// self-update
		cli.NewSelfUpdateCommand(ctx),

		// deploy
		cli.NewDeployCommand(ctx),

		// proxy commands
		proxy.NewRootCommand(ctx),

		// cloud commands
		cloud.NewRootCommand(ctx),
	)

	return nest
}

func main() {
	err := newNestCommand(&pkg.Context{}).Execute()
	if err != nil {
		global.LogI(global.LevelError, err.Error(), nil)

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
