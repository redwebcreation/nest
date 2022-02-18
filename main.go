package main

import (
	"fmt"
	"github.com/redwebcreation/nest/build"
	"github.com/redwebcreation/nest/cli"
	"github.com/redwebcreation/nest/cli/cloud"
	"github.com/redwebcreation/nest/cli/proxy"
	"github.com/redwebcreation/nest/pkg"
	logger2 "github.com/redwebcreation/nest/pkg/logger"
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
		Version:       fmt.Sprintf("%s, build %s", build.Version, build.Commit),
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx.Logger().Print(logger2.NewEvent(
				logger2.DebugLevel,
				"command invoked",
				logger2.Fields{
					"tag":     "command.invoke",
					"command": cmd.Name(),
				},
			))
		},
	}

	nest.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})

	// This flag is not actually used by any of the commands.
	// Its value is used in the init function in logger/server.go
	nest.PersistentFlags().StringP("config", "c", ctx.Home(), "set the logger config path")

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
	ctx, err := pkg.NewContext()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = newNestCommand(ctx).Execute()
	if err != nil {
		ctx.Logger().Print(logger2.NewEvent(logger2.ErrorLevel, err.Error(), nil))

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
