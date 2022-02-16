package cli

import (
	"fmt"
	"github.com/redwebcreation/nest/global"
	"github.com/spf13/cobra"
)

func runVersionCommand(cmd *cobra.Command, args []string) error {
	fmt.Printf("Nest version %s, build %s\n", global.Version, global.Commit)
	return nil
}

// NewVersionCommand creates a new `version` command.
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print nest's version",
		RunE:  runVersionCommand,
		Annotations: map[string]string{
			"medic":  "false",
			"config": "false",
		},
	}

	return cmd
}
