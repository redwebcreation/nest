package command

import (
	"fmt"

	"github.com/redwebcreation/nest/common"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
)

func runConfigCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("strategy:", common.ConfigReader.Strategy)
	fmt.Println("location:", common.ConfigReader.GetRepositoryLocation())
	fmt.Printf("current commit: %s\n", common.ConfigReader.LatestCommit[:7])

	configFiles, err := common.ConfigReader.Git.Files()
	if err != nil {
		return err
	}

	fmt.Println()
	util.NewTree(configFiles).Print(0)

	return nil
}

func ConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "prints the current configuration",
		RunE:  runConfigCommand,
	}

	return WithConfig(cmd)
}
