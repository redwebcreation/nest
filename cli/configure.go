package cli

import (
	"fmt"
	"regexp"

	"github.com/redwebcreation/nest/common"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
)

var strategy string
var provider string
var repository string

func ConfigureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "configure",
		Aliases: []string{
			"rcfg",
			"reconfigure",
		},
		Short: "update the global configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			usingFlags := strategy != "" || provider != "" || repository != ""

			if !usingFlags {
				common.ConfigReader.Strategy = util.Prompt("Choose a strategy", "remote", func(input string) bool {
					return input == "remote" || input == "local"
				})
				common.ConfigReader.Provider = util.Prompt("Choose a provider", "github", func(input string) bool {
					return input == "github" || input == "gitlab" || input == "bitbucket"
				})
				common.ConfigReader.Repository = util.Prompt("Enter a repository URL", common.ConfigReader.Repository, func(input string) bool {
					re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+")

					return re.MatchString(input)
				})
			} else {
				common.ConfigReader.Strategy = strategy
				common.ConfigReader.Provider = provider
				common.ConfigReader.Repository = repository
				err := common.ConfigReader.Validate()
				if err != nil {
					return err
				}

			}

			fmt.Println("\nSuccessfully configured the config resolver.")

			return common.ConfigReader.WriteOnDisk()
		},
	}

	cmd.Flags().StringVarP(&strategy, "strategy", "s", "", "strategy to use")
	cmd.Flags().StringVarP(&provider, "provider", "p", "", "provider to use")
	cmd.Flags().StringVarP(&repository, "repository", "r", "", "repository to use")

	return Decorate(cmd)
}
