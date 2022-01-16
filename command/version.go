package command

import (
	"fmt"

	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
)

func runVersionCommand(cmd *cobra.Command, args []string) error {
	_, _ = fmt.Fprintf(util.Stdout, "nest@%s\n", global.Version)
	return nil
}

// NewVersionCommand prints nest's version
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
