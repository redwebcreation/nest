package nest

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/redwebcreation/nest/config"
	"github.com/redwebcreation/nest/context"
	"regexp"

	"github.com/spf13/cobra"
)

type setupOptions struct {
	UsesFlags  bool
	Provider   string
	Repository string
	Branch     string
	Commit     string
}

func runSetupCommand(ctx *context.Context, opts *setupOptions) error {
	oldConfig, err := ctx.Config()
	hasConfig := err == nil

	if !opts.UsesFlags {
		// A default value for a select must be one of the options
		var repository = "github"

		if hasConfig {
			repository = oldConfig.Repository
		}

		prompt := &survey.Select{
			Message: "Select your provider:",
			Options: []string{"github", "gitlab", "bitbucket"},
			Default: repository,
		}
		err = survey.AskOne(prompt, &opts.Provider, survey.WithValidator(survey.Required), survey.WithStdio(ctx.In(), ctx.Out(), ctx.Err()))
		if err != nil {
			return err
		}
	}

	if !opts.UsesFlags {
		var repository string

		if hasConfig {
			repository = oldConfig.Repository
		}

		prompt := &survey.Input{
			Message: "Enter your repository:",
			Default: repository,
		}
		err = survey.AskOne(prompt, &opts.Repository, survey.WithValidator(func(ans any) error {
			re := regexp.MustCompile(`[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+`)
			if !re.MatchString(ans.(string)) {
				return fmt.Errorf("repository name must be alphanumeric and can contain hyphens and underscores")
			}
			return nil
		}), survey.WithStdio(ctx.In(), ctx.Out(), ctx.Err()))
		if err != nil {
			return err
		}
	}

	if !opts.UsesFlags {
		var branch string

		if hasConfig {
			branch = oldConfig.Branch
		}

		prompt := &survey.Input{
			Message: "Enter your branch:",
			Default: branch,
		}
		err = survey.AskOne(prompt, &opts.Branch, survey.WithValidator(survey.Required), survey.WithStdio(ctx.In(), ctx.Out(), ctx.Err()))
		if err != nil {
			return err
		}
	}

	config := config.Config{
		Provider:   opts.Provider,
		Repository: opts.Repository,
		Branch:     opts.Branch,
		Path:       ctx.ConfigFile(),
		StoreDir:   ctx.ConfigStoreDir(),
		Logger:     ctx.Logger(),
		Git: &config.Git{
			Logger: ctx.Logger(),
		},
	}
	if err = config.Save(); err != nil {
		return err
	}

	err = config.Clone()
	if err != nil {
		return err
	}

	fmt.Fprintln(ctx.Out(), "\nYou now need to run `nest use` to specify which version of the oldConfig you want to use.")

	return nil
}

// NewSetupCommand creates a new `setup` command.
func NewSetupCommand(ctx *context.Context) *cobra.Command {
	opts := &setupOptions{}

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "update the loggy config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.UsesFlags = cmd.Flags().NFlag() > 0
			return runSetupCommand(ctx, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Provider, "provider", "p", "github", "provider")
	cmd.Flags().StringVarP(&opts.Repository, "repository", "r", "", "repository")
	cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "main", "branch")

	return cmd
}
