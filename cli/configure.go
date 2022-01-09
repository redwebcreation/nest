package cli

import (
	"fmt"
	"github.com/me/nest/common"
	"github.com/me/nest/util"
	"github.com/spf13/cobra"
	"regexp"
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
				common.ConfigReader.ProviderURL = fmt.Sprintf("git@%s.com:", util.Prompt("Choose a provider", "github", func(input string) bool {
					return input == "github" || input == "gitlab" || input == "bitbucket"
				}))
				common.ConfigReader.Repository = util.Prompt("Enter a repository URL", common.ConfigReader.Repository, func(input string) bool {
					re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+")

					return re.MatchString(input)
				})
			} else {
				if strategy != "local" && strategy != "remote" {
					return fmt.Errorf("strategy must be either local or remote")
				}

				if provider != "github" && provider != "gitlab" && provider != "bitbucket" {
					return fmt.Errorf("provider must be either github, gitlab or bitbucket")
				}

				re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+")
				if !re.MatchString(repository) {
					return fmt.Errorf("repository must be in the format of <username>/<repository>")
				}

				common.ConfigReader.Strategy = strategy
				common.ConfigReader.ProviderURL = fmt.Sprintf("git@%s.com:", provider)
				common.ConfigReader.Repository = repository
			}

			fmt.Println("\nSuccessfully configured the config resolver.")

			return common.ConfigReader.WriteOnDisk()
		},
	}

	cmd.Flags().StringVarP(&strategy, "strategy", "s", "", "strategy to use")
	cmd.Flags().StringVarP(&provider, "provider", "p", "", "provider to use")
	cmd.Flags().StringVarP(&repository, "repository", "r", "", "repository to use")

	return cmd
}
