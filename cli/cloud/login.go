package cloud

import (
	"github.com/redwebcreation/nest/cloud"
	"github.com/redwebcreation/nest/util"
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
		util.Red.Render("Invalid token.")
	} else if err != nil {
		return err
	} else {
		util.Green.Render("Successfully logged in.")
	}

	return nil
}

// NewLoginCommand creates a new `login` command.
func NewLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login to nest cloud",
		Args:  cobra.ExactArgs(1),
		RunE:  runLoginCommand,
	}

	return cmd
}
