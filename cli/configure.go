package cli

import (
	"fmt"
	"regexp"

	"github.com/me/nest/global"
	"github.com/me/nest/ui"
	"github.com/spf13/cobra"
)

func ConfigureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "configure",
		Aliases: []string{
			"rcfg",
			"reconfigure",
		},
		Short: "Update the global configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			strategy, err := ui.Select{
				Question: "Choose a strategy: ",
				Choices:  []string{"remote"},
			}.Prompt()
			if err != nil {
				return err
			}
			global.ConfigLocatorConfig.Strategy = strategy

			provider, err := ui.Select{
				Question: "Choose a provider: ",
				Choices:  []string{"github", "gitlab", "bitbucket"},
			}.Prompt()
			if err != nil {
				return err
			}

			global.ConfigLocatorConfig.ProviderURL = fmt.Sprintf("git@%s.com:", provider)

			repository, err := ui.Input{
				Question: "Enter the repository URL: ",
				Validation: func(s string) bool {
					re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+")

					return re.MatchString(s)
				},
				Default: global.ConfigLocatorConfig.Repository,
			}.Prompt()
			if err != nil {
				return err
			}
			global.ConfigLocatorConfig.Repository = repository

			fmt.Println("\nSuccessfully configured the config resolver.")

			return global.ConfigLocatorConfig.SaveLocally()
		},
	}
	return cmd
}
