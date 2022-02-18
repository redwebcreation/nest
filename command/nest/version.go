package nest

import (
	"fmt"
	"github.com/redwebcreation/nest/build"
	"github.com/redwebcreation/nest/context"
	"github.com/spf13/cobra"
)

func runVersionCommand(ctx *context.Context) error {
	fmt.Fprintf(ctx.Out(), "Nest version %s, build %s\n\n", build.Version, build.Commit)
	return nil
}

// NewVersionCommand creates a new `version` command.
func NewVersionCommand(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print nest's version",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runVersionCommand(ctx)
		},
	}

	return cmd
}
