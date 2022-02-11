package command

import (
	"fmt"
	"os"

	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

// Configure sets defaults for the given command
func Configure(cmd *cobra.Command) {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	_, disableConfigLocator := cmd.Annotations["config"]
	_, disableMedic := cmd.Annotations["medic"]

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		global.InternalLogger.Log(
			global.LevelDebug,
			"command invoked",
			global.Fields{
				"tag":     "command.invoke",
				"command": cmd.Name(),
				"args":    args,
			},
		)

		if !disableConfigLocator {
			if _, err := os.Stat(global.LocatorConfigFile); err != nil {
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
	}

}
