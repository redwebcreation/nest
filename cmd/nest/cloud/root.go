package cloud

import (
	"github.com/redwebcreation/nest/context"
	"github.com/spf13/cobra"
)

// NewRootCommand creates a new `cloud` command.
func NewRootCommand(ctx *context.Context) *cobra.Command {
	root := &cobra.Command{
		Use:   "cloud",
		Short: "interact with nest cloud",
	}

	root.AddCommand(
		// login
		NewLoginCommand(ctx),
	)

	return root
}
