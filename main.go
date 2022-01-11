package main

import (
	"fmt"
	"github.com/redwebcreation/nest/command"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var commands = []*cobra.Command{
	command.DeployCommand(),
	command.MedicCommand(),
	command.ConfigCommand(),
	command.SelfUpdateCommand(),
	command.ConfigureCommand(),
	command.VersionCommand(),
}

func main() {
	nest := &cobra.Command{
		Use:   "nest",
		Short: "Service orchestrator",
		Long:  "Nest is a powerful service orchestrator for a single server.",
	}

	for _, cmd := range commands {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		nest.AddCommand(cmd)
	}

	nest.SetHelpCommand(&cobra.Command{
		Use:    "_help",
		Hidden: true,
	})
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
