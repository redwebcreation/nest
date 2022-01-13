package command

import (
	"encoding/json"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"os"
	"regexp"

	"github.com/redwebcreation/nest/common"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
)

var strategy string
var provider string
var repository string
var dir string

func runConfigureCommand(cmd *cobra.Command, args []string) error {
	if cmd.Flags().NFlag() == 0 {
		common.ConfigLocator.Strategy = util.Prompt("Choose a strategy", "remote", func(input string) bool {
			return input == "remote" || input == "local"
		})
		common.ConfigLocator.Provider = util.Prompt("Choose a provider", "github", func(input string) bool {
			return input == "github" || input == "gitlab" || input == "bitbucket"
		})
		common.ConfigLocator.Repository = util.Prompt("Enter a repository URL", common.ConfigLocator.Repository, func(input string) bool {
			re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+")

			return re.MatchString(input)
		})
	} else {
		if strategy != "" {
			common.ConfigLocator.Strategy = strategy
		}
		if provider != "" {
			common.ConfigLocator.Provider = provider
		}
		if repository != "" {
			common.ConfigLocator.Repository = repository
		}
		if dir != "" {
			common.ConfigLocator.Dir = dir
		}

		err := common.ConfigLocator.Validate()
		if err != nil {
			return err
		}
	}

	contents, err := json.Marshal(common.ConfigLocator)
	if err != nil {
		return err
	}

	err = os.WriteFile(global.ConfigLocatorConfigFile, contents, 0600)
	if err != nil {
		return err
	}

	fmt.Println("\nSuccessfully configured the config resolver.")
	return nil
}

// NewConfigureCommand configures the config locator
func NewConfigureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "configure",
		Aliases: []string{
			"rcfg",
			"reconfigure",
		},
		Short: "update the global configuration",
		RunE:  runConfigureCommand,
	}

	_ = LoadConfig()

	cmd.Flags().StringVarP(&strategy, "strategy", "s", common.ConfigLocator.Strategy, "strategy to use")
	cmd.Flags().StringVarP(&provider, "provider", "p", common.ConfigLocator.Provider, "provider to use")
	cmd.Flags().StringVarP(&repository, "repository", "r", common.ConfigLocator.Repository, "repository to use")
	cmd.Flags().StringVarP(&dir, "dir", "d", common.ConfigLocator.Dir, "dir in repo to use as root")

	return cmd
}
