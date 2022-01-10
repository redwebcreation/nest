package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/redwebcreation/nest/common"
	"github.com/redwebcreation/nest/global"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/redwebcreation/nest/cli"
	"github.com/redwebcreation/nest/cli/proxy"
)

var commands = []*cobra.Command{
	proxy.RootCommand(),
	cli.DeployCommand(),
	cli.MedicCommand(),
	cli.ConfigCommand(),
	cli.SelfUpdateCommand(),
	cli.ConfigureCommand(),
	cli.VersionCommand(),
}

func main() {
	nest := &cobra.Command{
		Use:   "nest",
		Short: "Service orchestrator",
		Long:  "Nest is a powerful service orchestrator for a single server.",
	}

	for _, command := range commands {
		command.SilenceUsage = true
		command.SilenceErrors = true

		nest.AddCommand(command)
	}

	// hide the help command
	nest.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})

	nest.PersistentPreRunE = prerun
	nest.SilenceErrors = true

	err := nest.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error: "+err.Error())
		os.Exit(1)
	}
}

func init() {
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Fprintf(os.Stderr, "error: git is not installed")
		os.Exit(1)
	}
}

func prerun(cmd *cobra.Command, _ []string) error {
	commandName := cmd.Name()

	if _, err := os.Stat(global.ConfigLocatorConfigFile); err != nil {
		if commandName == "configure" {
			return nil
		}

		return fmt.Errorf("run `nest configure` to setup nest")
	}

	reader, err := common.LoadConfigReader()
	if err != nil && commandName == "configure" {
		common.ConfigReader = common.NewConfigReader()
		return nil
	} else if err != nil {
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
