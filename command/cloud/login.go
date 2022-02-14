package cloud

import (
	"github.com/redwebcreation/nest/cloud"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
)

func runLoginCommand(cmd *cobra.Command, args []string) error {
	token := args[0]

	err := cloud.SetToken(token)
	if err != nil {
		return err
	}

	err = cloud.NewClient(token).Ping()
	if err == cloud.ErrResourceNotFound {
		pkg.Red.Render("Invalid token.")
	} else if err != nil {
		return err
	} else {
		pkg.Green.Render("Successfully logged in.")
	}

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
