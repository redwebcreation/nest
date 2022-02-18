package cloud

import (
	"fmt"
	"github.com/redwebcreation/nest/cloud"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"os"
)

type loginOptions struct {
	token string
}

func runLoginCommand(ctx *pkg.Context, opts *loginOptions) error {
	err := os.WriteFile(ctx.CloudTokenFile(), []byte(opts.token), 0600)
	if err != nil {
		return err
	}

	err = cloud.NewClient(opts.token).Ping()
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
func NewLoginCommand(ctx *pkg.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login to nest cloud",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoginCommand(ctx, &loginOptions{
				token: args[0],
			})
		},
	}

	return cmd
}
