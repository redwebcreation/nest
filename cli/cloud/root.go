package cloud

import (
	"github.com/redwebcreation/nest/cloud"
	"github.com/spf13/cobra"
)

var commands = []*cobra.Command{
	NewLoginCommand(),
}

// NewRootCommand creates a new `cloud` command.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "cloud",
		Short: "interact with nest cloud",
	}

	root.PersistentFlags().StringVar(&cloud.Endpoint, "endpoint", cloud.Endpoint, "nest cloud endpoint")

	for _, command := range commands {
		root.AddCommand(command)
	}

	return root
}
