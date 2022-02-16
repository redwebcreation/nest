package proxy

import (
	"github.com/spf13/cobra"
)

// NewRootCommand returns a new instance of the proxy root command
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "proxy",
		Short: "manage the proxy",
	}

	for _, cmd := range []*cobra.Command{NewRunCommand()} {
		root.AddCommand(cmd)
	}

	return root
}
