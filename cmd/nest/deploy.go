package nest

import (
	"fmt"
	"github.com/redwebcreation/nest/context"
	"github.com/redwebcreation/nest/deploy"
	"github.com/spf13/cobra"
	"io"
)

func runDeployCommand(ctx *context.Context) error {
	config, err := ctx.ServerConfig()
	if err != nil {
		return err
	}

	deployment := deploy.NewDeployment(config, ctx.ManifestManager())

	go func() {
		err = deployment.Start()
		if err != nil {
			deployment.Events <- deploy.Event{
				Service: nil,
				Value:   deploy.ErrDeploymentFailed,
			}
		}
	}()

	for event := range deployment.Events {
		if event.Value == deploy.ErrDeploymentFailed {
			fmt.Fprintln(ctx.Out(), "Deployment failed")
			break
		}

		if event.Value == io.EOF {
			break
		}

		if event.Service != nil {
			fmt.Fprintf(ctx.Out(), "%s: %v\n", event.Service.Name, event.Value)
		} else {
			fmt.Fprintf(ctx.Out(), "loggy: %v\n", event.Value)
		}
	}

	return ctx.ManifestManager().Save(deployment.Manifest)
}

// NewDeployCommand creates a new `deploy` command.
func NewDeployCommand(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy the config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeployCommand(ctx)
		},
	}

	return cmd
}
