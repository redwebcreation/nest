package cli

import (
	"fmt"
	"github.com/me/nest/global"
	"github.com/me/nest/ui"
	"github.com/spf13/cobra"
	"regexp"
)

func ReconfigureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "reconfigure",
		Aliases: []string{
			"rcfg",
		},
		Short: "Update the current global configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			strategy, err := ui.Select{
				Question: "Choose a strategy: ",
				Choices:  []string{"remote"},
			}.Prompt()
			if err != nil {
				return err
			}
			global.ConfigResolver.Strategy = strategy

			provider, err := ui.Select{
				Question: "Choose a provider: ",
				Choices:  []string{"GitHub", "GitLab", "Bitbucket"},
			}.Prompt()
			if err != nil {
				return err
			}
			global.ConfigResolver.Provider = provider

			transportMode, err := ui.Select{
				Question: "Choose a transport mode: ",
				Choices:  []string{"ssh", "https"},
			}.Prompt()
			if err != nil {
				return err
			}
			global.ConfigResolver.TransportMode = transportMode

			repository, err := ui.Input{
				Question: "Enter the repository URL: ",
				Validation: func(s string) bool {
					re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+")

					return re.MatchString(s)
				},
				Default: global.ConfigResolver.Repository,
			}.Prompt()
			if err != nil {
				return err
			}
			global.ConfigResolver.Repository = repository

			fmt.Println("\nSuccessfully configured the finder.")

			return global.ConfigResolver.Write()
		},
	}
	return cmd
}
