package main

import (
	"fmt"
	"github.com/redwebcreation/nest/command"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var commands = []*cobra.Command{
	command.NewDeployCommand(),
	command.NewMedicCommand(),
	command.NewConfigCommand(),
}

var standalone = []*cobra.Command{
	command.NewConfigureCommand(),
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
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		command.WithConfig(cmd)
		nest.AddCommand(cmd)
	}

	for _, cmd := range standalone {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		nest.AddCommand(cmd)
	}

	err := nest.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error: "+err.Error())
		os.Exit(1)
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
