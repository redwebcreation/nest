package config

import (
	"fmt"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

func runUseCommand(cmd *cobra.Command, args []string) error {
	fmt.Println()

	commits, err := pkg.Git.ListCommits(pkg.Locator.ConfigPath(), pkg.Locator.Branch)
	if err != nil {
		return err
	}

	pkg.Gray.Render("Inspecting " + strconv.Itoa(len(commits)) + " commits...")

	var commit string

	if len(args) > 0 {
		var possibilities []pkg.Commit

		for _, c := range commits {
			if strings.HasPrefix(c.Hash, args[0]) {
				possibilities = append(possibilities, c)
			}
		}

		if len(possibilities) != 1 {
			// todo:
			//util.PrintE(util.Red.Fg()+"\n  Could not find a unique match for %s.\n"+util.Reset()+util.Gray.Fg(), args[0])
			//util.PrintE()
			//util.PrintE("  Candidates")
			//
			//for _, possibility := range possibilities {
			//	util.PrintE("  * %s %s\n", possibility.Hash[:7], possibility.Message)
			//}
			//
			//util.FatalE(util.Reset())
			os.Exit(0)
		} else {
			//fmt.Printf(util.Gray.Fg()+"  Found unique commit %s '%s'.\n\n"+util.Reset(), possibilities[0].Hash[:7], possibilities[0].Message)
			commit = possibilities[0].Hash
		}
	}

	err = pkg.Locator.LoadCommit(commit)
	if err != nil {
		return err
	}

	// Using pkg.Locator.Commit rather than commit to get the real commit used if no arguments were passed
	pkg.Green.Render("  Updated the locator config. Now using " + pkg.Locator.Commit[:7] + ".")

	return nil
}

// NewUseCommand sets the command to use for the config locator
func NewUseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use [commit]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Use a specific commit",
		RunE:  runUseCommand,
	}

	return cmd
}
