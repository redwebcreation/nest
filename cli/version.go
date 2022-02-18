package cli

import (
	"fmt"
	"github.com/redwebcreation/nest/build"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

func runVersionCommand(ctx *pkg.Context) error {
	fmt.Fprintf(ctx.Out(), "Nest version %s, build %s\n\n", build.Version, build.Commit)
	return nil
}

// NewVersionCommand creates a new `version` command.
func NewVersionCommand(ctx *pkg.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print nest's version",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runVersionCommand(ctx)
		},
	}

	return cmd
}
