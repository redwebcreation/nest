package cloud

import (
	"fmt"
	"github.com/redwebcreation/nest/cloud"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
)

func runLoginCommand(cmd *cobra.Command, args []string) error {
	token := args[0]

	err := cloud.AddToken(token)
	if err != nil {
		return err
	}

	fmt.Println(util.Green.Fg() + "Successfully logged in." + util.Reset())

	return nil
}

func NewLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login to nest cloud",
		Args:  cobra.ExactArgs(1),
		RunE:  runLoginCommand,
	}

	return cmd
}
