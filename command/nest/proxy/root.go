package proxy

import (
	"github.com/redwebcreation/nest/context"
	"github.com/spf13/cobra"
)

// NewRootCommand returns a new instance of the proxy root command
func NewRootCommand(ctx *context.Context) *cobra.Command {
	root := &cobra.Command{
		Use:   "proxy",
		Short: "manage the proxy",
	}

	root.AddCommand(
		// run
		NewRunCommand(ctx),
	)

	return root
}
