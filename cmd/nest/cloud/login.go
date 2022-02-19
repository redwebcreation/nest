package cloud

import (
	"fmt"
	"github.com/redwebcreation/nest/cloud"
	"github.com/redwebcreation/nest/context"
	"github.com/spf13/cobra"
	"os"
)

type loginOptions struct {
	id          string
	accessToken string
}

func runLoginCommand(ctx *context.Context, opts *loginOptions) error {
	err := os.WriteFile(ctx.CloudCredentialsFile(), []byte(fmt.Sprintf("%s:%s", opts.id, opts.accessToken)), 0600)
	if err != nil {
		return err
	}

	client, err := ctx.CloudClient()
	if err != nil {
		return err
	}

	err = client.Ping()
	if err == cloud.ErrResourceNotFound {
		fmt.Fprintln(ctx.Out(), "Invalid token.")
	} else if err != nil {
		return err
	} else {
		fmt.Fprintln(ctx.Out(), "Successfully logged in.")
	}

	return nil
}

// NewLoginCommand creates a new `login` command.
func NewLoginCommand(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login to nest cloud",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args[0]) != 45 {
				return fmt.Errorf("invalid token")
			}

			return runLoginCommand(ctx, &loginOptions{
				id:          args[0][:22],
				accessToken: args[0][23:],
			})
		},
	}

	return cmd
}
