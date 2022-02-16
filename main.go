package main

import (
	"fmt"
	"github.com/redwebcreation/nest/cli"
	"github.com/redwebcreation/nest/cli/cloud"
	"github.com/redwebcreation/nest/cli/config"
	"github.com/redwebcreation/nest/cli/proxy"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"os"
)

var commands = []*cobra.Command{
	cloud.NewRootCommand(),
	proxy.NewRootCommand(),
	config.NewRootCommand(),
	cli.NewDeployCommand(),
	cli.NewMedicCommand(),
	cli.NewSetupCommand(),
	cli.NewVersionCommand(),
	cli.NewSelfUpdateCommand(),
}

func newNestCommand() *cobra.Command {
	nest := &cobra.Command{
		Use:     "nest",
		Short:   "Service orchestrator",
		Long:    "Nest is a powerful service orchestrator for a single server.",
		Version: fmt.Sprintf("%s, build %s", global.Version, global.Commit),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			_, disableConfigLocator := cmd.Annotations["config"]
			_, disableMedic := cmd.Annotations["medic"]

			global.LogI(
				global.LevelDebug,
				"command invoked",
				global.Fields{
					"tag":     "command.invoke",
					"command": cmd.Name(),
				},
			)

			if !disableConfigLocator {
				if _, err := os.Stat(global.GetLocatorConfigFile()); err != nil {
					return fmt.Errorf("run `nest setup` to setup nest")
				}

				err := pkg.Locator.Load()
				if err != nil {
					return err
				}
			}

			if disableMedic {
				return nil
			}

			return pkg.DiagnoseConfiguration().MustPass()
		},
	}

	nest.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})

	nest.CompletionOptions.DisableDefaultCmd = true

	// This flag is not actually used by any of the commands.
	// Its value is used in the init function in global/config.go
	nest.PersistentFlags().StringP("config", "c", global.ConfigHome, "set the global config path")

	for _, command := range commands {
		nest.AddCommand(command)
	}

	return nest
}

func main() {
	cmd := newNestCommand()

	err := cmd.Execute()
	if err != nil {
		global.LogI(global.LevelError, err.Error(), nil)

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
