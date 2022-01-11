package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

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

	nest.AddCommand(commands...)
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
