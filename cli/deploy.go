package cli

import (
	"fmt"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"io"
)

func runDeployCommand(ctx *pkg.Context) error {
	config, err := ctx.ServerConfig()
	if err != nil {
		return err
	}

	deployment := pkg.NewDeployment(config, ctx.ManifestManager())

	go func() {
		err = deployment.Start()
		if err != nil {
			deployment.Events <- pkg.Event{
				Service: nil,
				Value:   pkg.ErrDeploymentFailed,
			}
		}
	}()

	for event := range deployment.Events {
		if event.Value == pkg.ErrDeploymentFailed {
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

	return deployment.Manifest.Save(ctx.ManifestFile(deployment.ID))
}

// NewDeployCommand creates a new `deploy` command.
func NewDeployCommand(ctx *pkg.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy the config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeployCommand(ctx)
		},
	}

	return cmd
}
