package main

import (
	"fmt"
	"github.com/me/nest/common"
	"github.com/me/nest/global"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"

	"github.com/me/nest/cli"
	"github.com/me/nest/cli/proxy"
	"github.com/spf13/cobra"
)

func main() {
	// check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("Git is not installed. Please install git and try again.")
		os.Exit(1)
	}

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

		if _, err := os.Stat(global.ConfigLocatorConfigFile); err != nil {
			if commandName == "configure" {
				return nil
			}

			return fmt.Errorf("run `nest configure` to setup nest")
		}

		reader, err := common.LoadConfigReader()
		if err != nil {
			return err
		}

		common.ConfigReader = reader

		contents, err := reader.Read("nest.yaml")
		if err != nil {
			return err
		}

		var config common.Configuration

		err = yaml.Unmarshal(contents, &config)
		if err != nil {
			return err
		}

		common.Config = &config

		if commandName == "medic" {
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
