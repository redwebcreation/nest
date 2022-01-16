package command

import (
	"fmt"

	"github.com/redwebcreation/nest/pkg"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
)

func runConfigCommand(cmd *cobra.Command, args []string) error {
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
	util.NewTree(configFiles).Print(0)

	return nil
}

// NewConfigCommand prints the current configuration for the config locator
func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "prints the current configuration",
		RunE:  runConfigCommand,
		Annotations: map[string]string{
			"medic": "false",
		},
	}

	return cmd
}
