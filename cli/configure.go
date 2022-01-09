package cli

import (
	"github.com/spf13/cobra"
)

func ConfigureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "configure",
		Aliases: []string{
			"rcfg",
			"reconfigure",
		},
		Short: "update the global configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
			//strategy, err := ui.Select{
			//	Question: "Choose a strategy: ",
			//	Choices:  []string{"remote"},
			//}.Prompt()
			//if err != nil {
			//	return err
			//}
			//common.ConfigLocatorConfig.Strategy = strategy
			//
			//provider, err := ui.Select{
			//	Question: "Choose a provider: ",
			//	Choices:  []string{"github", "gitlab", "bitbucket"},
			//}.Prompt()
			//if err != nil {
			//	return err
			//}
			//
			//common.ConfigLocatorConfig.ProviderURL = fmt.Sprintf("git@%s.com:", provider)
			//
			//repository, err := ui.Input{
			//	Question: "Enter the repository URL: ",
			//	Validation: func(s string) bool {
			//		re := regexp.MustCompile("[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+")
			//
			//		return re.MatchString(s)
			//	},
			//	Default: common.ConfigLocatorConfig.Repository,
			//}.Prompt()
			//if err != nil {
			//	return err
			//}
			//common.ConfigLocatorConfig.Repository = repository
			//
			//fmt.Println("\nSuccessfully configured the config resolver.")
			//
			//return common.ConfigLocatorConfig.SaveLocally()
		},
	}
	return cmd
}
