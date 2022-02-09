package proxy

import (
	"github.com/redwebcreation/nest/command"
	"github.com/spf13/cobra"
)

// NewRootProxyCommand returns a new instance of the proxy root command
func NewRootProxyCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "proxy",
		Short: "manage the proxy",
	}

	for _, cmd := range []*cobra.Command{NewRunCommand()} {
		command.Configure(cmd)

		root.AddCommand(cmd)
	}

	return root
}
