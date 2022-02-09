package config

import (
	"fmt"

	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

func runPullCommand(cmd *cobra.Command, args []string) error {
	out, err := pkg.Locator.VCS.Pull(pkg.Locator.ConfigPath(), pkg.Locator.Branch)

	fmt.Println(string(out))

	return err
}

func NewPullCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pulls the latest version of the configuration",
		RunE:  runPullCommand,
	}

	return cmd
}
