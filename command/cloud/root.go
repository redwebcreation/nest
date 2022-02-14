package cloud

import (
	"github.com/redwebcreation/nest/cloud"
	"github.com/redwebcreation/nest/command"
	"github.com/spf13/cobra"
)

var commands = []*cobra.Command{
	NewLoginCommand(),
}

func NewRootConfigCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "cloud",
		Short: "interact with nest cloud",
	}

	root.PersistentFlags().StringVar(&cloud.Endpoint, "endpoint", cloud.Endpoint, "nest cloud endpoint")

	for _, cmd := range commands {
		command.Configure(cmd)

		root.AddCommand(cmd)
	}

	return root
}
