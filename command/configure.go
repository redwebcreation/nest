package command

import (
	"encoding/json"
	"fmt"
	"github.com/redwebcreation/nest/global"
	"os"
	"regexp"

	"github.com/redwebcreation/nest/pkg"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
)

var strategy string
var provider string
var repository string
var branch string
var dir string

func runConfigureCommand(cmd *cobra.Command, args []string) error {
	if cmd.Flags().NFlag() == 0 {
		pkg.Config.Strategy = util.Prompt("Choose a strategy", "remote", func(input string) bool {
			return input == "remote" || input == "local"
		})
		pkg.Config.Provider = util.Prompt("Choose a provider", "github", func(input string) bool {
			return input == "github" || input == "gitlab" || input == "bitbucket"
		})
		pkg.Config.Repository = util.Prompt("Enter a repository URL", pkg.Config.Repository, func(input string) bool {
			re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+")

			return re.MatchString(input)
		})
		pkg.Config.Branch = util.Prompt("Enter a branch", pkg.Config.Branch, func(input string) bool {
			return input != ""
		})
	} else {
		if strategy != "" {
			pkg.Config.Strategy = strategy
		}
		if provider != "" {
			pkg.Config.Provider = provider
		}
		if repository != "" {
			pkg.Config.Repository = repository
		}
		if dir != "" {
			pkg.Config.Dir = dir
		}
		if branch != "" {
			pkg.Config.Branch = branch
		}

		err := pkg.Config.Validate()
		if err != nil {
			return err
		}
	}

	contents, err := json.Marshal(pkg.Config)
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

	cmd.Flags().StringVarP(&strategy, "strategy", "s", pkg.Config.Strategy, "strategy to use")
	cmd.Flags().StringVarP(&provider, "provider", "p", pkg.Config.Provider, "provider to use")
	cmd.Flags().StringVarP(&repository, "repository", "r", pkg.Config.Repository, "repository to use")
	cmd.Flags().StringVarP(&branch, "branch", "b", pkg.Config.Branch, "branch to use")
	cmd.Flags().StringVarP(&dir, "dir", "d", pkg.Config.Dir, "dir in repo to use as root")

	return cmd
}
