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

func runSetupCommand(cmd *cobra.Command, args []string) error {
	pkg.Locator.Provider = util.Prompt("Choose a provider", pkg.Locator.Provider, func(input string) bool {
		return input == "github" || input == "gitlab" || input == "bitbucket"
	})

	pkg.Locator.Repository = util.Prompt("Enter a repository URL", pkg.Locator.Repository, func(input string) bool {
		re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+")

		return re.MatchString(input)
	})

	pkg.Locator.Branch = util.Prompt("Enter a branch", pkg.Locator.Branch, func(input string) bool {
		return input != ""
	})

	err := pkg.Locator.CloneConfig()
	if err != nil {
		util.PrintE(err)

		return cmd.Execute()
	}

	commits, err := pkg.Locator.VCS.ListCommits(pkg.Locator.ConfigPath(), pkg.Locator.Branch)
	if err != nil {
		return err
	}

	pkg.Locator.Commit = util.Prompt("Enter a full commit hash", commits[0].Hash, func(input string) bool {
		for _, commit := range commits {
			if commit.Hash == input { // todo: support for short commit hashes, re-use logic from `nest config use`
				return true
			}
		}

		return false
	})

	contents, err := json.Marshal(pkg.Locator)
	if err != nil {
		return err
	}

	err = os.WriteFile(global.LocatorConfigFile, contents, 0600)
	if err != nil {
		return err
	}

	fmt.Println("\nSuccessfully configured the config resolver.")
	return nil
}

// NewSetupCommand configures the config locator
func NewSetupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "update the global configuration",
		RunE:  runSetupCommand,
		Annotations: map[string]string{
			"medic":  "false",
			"config": "false",
		},
	}

	// load defaults from the config file
	// it's okay if it fails, we're reconfiguring it anyway
	_ = pkg.Locator.Load()

	return cmd
}
