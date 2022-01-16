package config

import (
	"fmt"

	"github.com/redwebcreation/nest/pkg"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
)

func runShowCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("strategy:", pkg.Config.Strategy)
	fmt.Println("location:", pkg.Config.GetRepositoryLocation())
	fmt.Printf("current commit: %s\n", pkg.Config.Commit[:7])
	fmt.Println("branch:", pkg.Config.Branch)
	if pkg.Config.Dir != "" {
		fmt.Println("subdir:", pkg.Config.Dir)
	}

	configFiles, err := pkg.Config.Git.Tree()
	if err != nil {
		return err
	}

	fmt.Println()
	util.PrintTree(configFiles)

	return nil
}

// NewShowCommand prints the current configuration for the config locator
func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "prints the current configuration",
		RunE:  runShowCommand,
	}

	return cmd
}
