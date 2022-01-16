package main

import (
	"fmt"
	"github.com/redwebcreation/nest/command"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var commands = []*cobra.Command{
	command.NewDeployCommand(),
	command.NewMedicCommand(),
	command.NewConfigCommand(),
	command.NewSetupCommand(),
	command.NewVersionCommand(),
	command.NewSelfUpdateCommand(),
}

var nest = &cobra.Command{
	Use:   "nest",
	Short: "Service orchestrator",
	Long:  "Nest is a powerful service orchestrator for a single server.",
}

func main() {
	for _, cmd := range commands {
		configure(cmd)
		nest.AddCommand(cmd)
	}

	err := nest.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error: "+err.Error())
		os.Exit(1)
	}
}

func configure(cmd *cobra.Command) {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	_, disableConfigLocator := cmd.Annotations["config"]
	_, disableMedic := cmd.Annotations["medic"]

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if !disableConfigLocator {
			if _, err := os.Stat(global.ConfigLocatorConfigFile); err != nil {
				return fmt.Errorf("run `nest setup` to setup nest")
			}

			if err := pkg.LoadConfig(); err != nil {
				return err
			}
		}

		if disableMedic {
			return nil
		}

		return pkg.DiagnoseConfiguration().MustPass()
	}
}

func init() {
	if _, err := exec.LookPath("git"); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: git is not installed")
		os.Exit(1)
	}

	nest.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})
}
