package config

import (
	"bytes"
	"fmt"
	"github.com/redwebcreation/nest/pkg"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
	"strings"
)

func runUseCommand(cmd *cobra.Command, args []string) error {
	fmt.Println()

	commits, err := pkg.Locator.VCS.ListCommits(pkg.Locator.ConfigPath(), pkg.Locator.Branch)
	if err != nil {
		return err
	}

	fmt.Printf(util.Gray.Fg()+"  Found %d candidates.\n"+util.Reset(), len(commits))

	var commit string

	if len(args) > 0 {
		var possibilities []util.Commit

		for _, c := range commits {
			if strings.HasPrefix(c.Hash, args[0]) {
				possibilities = append(possibilities, c)
			}
		}

		if len(possibilities) != 1 {
			var buffer bytes.Buffer

			buffer.WriteString(fmt.Sprintf(util.Red.Fg()+"\n  Could not find a unique match for %s.\n"+util.Reset()+util.Gray.Fg(), args[0]))
			buffer.WriteString("\n")
			buffer.WriteString("  Candidates:\n")

			for _, possibility := range possibilities {
				buffer.WriteString(fmt.Sprintf("  * %s %s\n", possibility.Hash[:7], possibility.Message))
			}

			buffer.WriteString(util.Reset())

			util.Fatal(buffer.String())
		} else {
			fmt.Printf(util.Gray.Fg()+"  Found %s.\n\n"+util.Reset(), possibilities[0])
			commit = possibilities[0].Hash
		}
	}

	err = pkg.Locator.LoadCommit(commit)
	if err != nil {
		return err
	}

	// Using pkg.Locator.Commit rather than commit to get the real commit used if no arguments were passed
	fmt.Printf(util.Green.Fg()+"  Updated the locator config. Now using %s.\n"+util.Reset(), pkg.Locator.Commit[:7])

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
