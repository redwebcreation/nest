package cli

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
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
			fmt.Printf("current commit (%s): %s\n", common.ConfigReader.Head.Hash.String()[:7], common.ConfigReader.Head.Message)

			files, err := common.ConfigReader.Head.Files()
			if err != nil {
				return err
			}

			var allFiles []string

			_ = files.ForEach(func(file *object.File) error {
				allFiles = append(allFiles, file.Name)

				return nil
			})

			fmt.Println()
			util.NewTree(allFiles).Print(0)

			return nil
		},
	}
}
