package cli

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

type useOptions struct {
	commit string
}

func runUseCommand(ctx *pkg.Context, opts *useOptions) error {
	config, err := ctx.Config()
	if err != nil {
		return err
	}

	err = config.Pull()
	if err != nil {
		return err
	}

	commits, err := config.Git.ListCommits(config.StorePath(), config.Branch)
	if err != nil {
		return err
	}

	fmt.Fprintf(ctx.Out(), "Inspecting %d commits...\n", len(commits))

	if opts.commit == "" {
		prompt := survey.Select{
			Message: "Select a commit to use",
			Options: commits.Hashes(),
		}
		err = survey.AskOne(&prompt, &opts.commit, survey.WithStdio(ctx.In(), ctx.Out(), ctx.Err()))
		if err != nil {
			return err
		}
	}

	if len(opts.commit) != 40 {
		return fmt.Errorf("invalid commit hash (must be full): %s", opts.commit)
	}

	err = config.LoadCommit(opts.commit)
	if err != nil {
		return err
	}

	fmt.Fprintf(ctx.Out(), "Updated the config. Now using %s.\n", opts.commit[:7])

	return nil
}

// NewUseCommand creates a new `use` command
func NewUseCommand(ctx *pkg.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use [commit]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Use a specific commit",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := &useOptions{
				commit: "",
			}

			if len(args) > 0 {
				opts.commit = args[0]
			}

			return runUseCommand(ctx, opts)
		},
	}

	return cmd
}
