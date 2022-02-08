package config

import (
	"fmt"

	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

func runPullCommand(cmd *cobra.Command, args []string) error {
	git, err := pkg.Config.LocalClone()
	if err != nil {
		return err
	}

	out, err := git.Pull(pkg.Config.Branch)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}

func NewPullCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pulls the latest version of the configuration",
		RunE:  runPullCommand,
	}

	return cmd
}
