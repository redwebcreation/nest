package command

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/redwebcreation/nest/global"
	"os"
	"regexp"

	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

var RepositoryNameValidator = func(ans any) error {
	re := regexp.MustCompile(`[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+`)
	if !re.MatchString(ans.(string)) {
		return fmt.Errorf("repository name must be alphanumeric and can contain hyphens and underscores")
	}
	return nil
}

func runSetupCommand(cmd *cobra.Command, args []string) error {
	var setup = []*survey.Question{
		{
			Name: "provider",
			Prompt: &survey.Select{
				Message: "Select your provider:",
				Options: []string{"github", "gitlab", "bitbucket"},
				Default: pkg.Locator.Provider,
			},
			Validate: survey.Required,
		},
		{
			Name: "repository",
			Prompt: &survey.Input{
				Message: "Enter your repository:",
				Default: pkg.Locator.Repository,
			},
			Validate: RepositoryNameValidator,
		},
	}
	err := survey.Ask(setup, &pkg.Locator)
	if err != nil {
		return err
	}

	err = pkg.Locator.CloneConfig()
	if err != nil {
		pkg.Stderr.Print(err)

		return cmd.Execute()
	}

	commits, err := pkg.Git.ListCommits(pkg.Locator.ConfigPath(), pkg.Locator.Branch)
	if err != nil {
		return err
	}

	promptCommit := &survey.Select{
		Message: "Select your commit:",
		Options: commits.Hashes(),
		Default: pkg.Locator.Commit,
	}
	err = survey.AskOne(promptCommit, &pkg.Locator.Commit)
	if err != nil {
		return err
	}

	contents, err := json.Marshal(pkg.Locator)
	if err != nil {
		return err
	}

	err = os.WriteFile(global.GetLocatorConfigFile(), contents, 0600)
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
