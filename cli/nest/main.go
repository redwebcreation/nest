package main

import (
	"fmt"
	"github.com/redwebcreation/nest/build"
	"github.com/redwebcreation/nest/command/nest"
	"github.com/redwebcreation/nest/command/nest/cloud"
	"github.com/redwebcreation/nest/command/nest/proxy"
	"github.com/redwebcreation/nest/context"
	"github.com/redwebcreation/nest/loggy"
	"github.com/spf13/cobra"
	"os"
)

func newNestCommand(ctx *context.Context) *cobra.Command {
	cli := &cobra.Command{
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
			ctx.Logger().Print(loggy.NewEvent(
				loggy.DebugLevel,
				"command invoked",
				loggy.Fields{
					"tag":     "command.invoke",
					"command": cmd.Name(),
				},
			))
		},
	}

	cli.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})

	// This flag is not actually used by any of the commands.
	// Its value is used in the init function in loggy/server.go
	cli.PersistentFlags().StringP("config", "c", ctx.Home(), "set the loggy config path")

	cli.AddCommand(
		// version
		nest.NewVersionCommand(ctx),

		// setup
		nest.NewSetupCommand(ctx),

		// use
		nest.NewUseCommand(ctx),

		// medic
		nest.NewMedicCommand(ctx),

		// self-update
		nest.NewSelfUpdateCommand(ctx),

		// deploy
		nest.NewDeployCommand(ctx),

		// proxy commands
		proxy.NewRootCommand(ctx),

		// cloud commands
		cloud.NewRootCommand(ctx),
	)

	return cli
}

func main() {
	ctx, err := context.NewContext()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = newNestCommand(ctx).Execute()
	if err != nil {
		ctx.Logger().Print(loggy.NewEvent(loggy.ErrorLevel, err.Error(), nil))

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
