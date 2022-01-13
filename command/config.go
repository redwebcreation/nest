package command

import (
	"fmt"

	"github.com/redwebcreation/nest/common"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
)

func runConfigCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("strategy:", common.ConfigLocator.Strategy)
	fmt.Println("location:", common.ConfigLocator.GetRepositoryLocation())
	fmt.Printf("current commit: %s\n", common.ConfigLocator.Commit[:7])

	configFiles, err := common.ConfigLocator.Git.Tree()
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
	}

	return cmd
}
