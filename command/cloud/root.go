package cloud

import "github.com/spf13/cobra"

func NewRootConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud",
		Short: "interact with nest cloud",
	}

	cmd.AddCommand(NewLoginCommand())

	return cmd
}
