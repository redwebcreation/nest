package cli

import (
	"fmt"
	"github.com/me/nest/common"
	"github.com/me/nest/util"
	"github.com/spf13/cobra"
)

func ConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "prints the current configuration",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println("strategy:", common.ConfigReader.Strategy)
			fmt.Println("location:", common.ConfigReader.ProviderURL+common.ConfigReader.Repository)
			fmt.Printf("current commit: %s\n", common.ConfigReader.LatestCommit[:7])

			configFiles, err := common.ConfigReader.Git.Files()
			if err != nil {
				return err
			}

			fmt.Println()
			util.NewTree(configFiles).Print(0)

			return nil
		},
	}
}
