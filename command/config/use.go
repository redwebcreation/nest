package config

import (
	"fmt"
	"strings"

	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

func runUseCommand(cmd *cobra.Command, args []string) error {
	commits, err := pkg.Config.Git.Commits()
	if err != nil {
		return err
	}

	var commit string

	if args[0] == "latest" {
		commit = commits[0]
	} else {
		for _, c := range commits {
			if strings.HasPrefix(c, args[0]) {
				commit = c
				break
			}
		}

		if commit == "" {
			return fmt.Errorf("commit not found")
		}
	}

	err = pkg.LoadConfigFromCommit(commit)
	if err != nil {
		return err
	}

	fmt.Printf("Using %s.\n", commit)

	return nil
}

func NewUseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use [commit]",
		Args:  cobra.ExactArgs(1),
		Short: "Use a specific commit",
		RunE:  runUseCommand,
	}

	return cmd
}
