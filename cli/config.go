package cli

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/me/nest/global"
	"github.com/me/nest/util"
	"github.com/spf13/cobra"
)

func ConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "prints the current configuration",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println("strategy:", global.ConfigLocatorConfig.Strategy)
			fmt.Println("location:", global.ConfigLocatorConfig.ProviderURL+global.ConfigLocatorConfig.Repository)
			fmt.Printf("current commit (%s): %s\n", global.ConfigLocatorConfig.Head.Hash.String()[:7], global.ConfigLocatorConfig.Head.Message)

			files, err := global.ConfigLocatorConfig.Head.Files()
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
