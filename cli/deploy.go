package cli

import (
	"fmt"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"io"
)

func runDeployCommand(cmd *cobra.Command, args []string) error {
	config, err := pkg.Locator.Resolve()
	if err != nil {
		return err
	}

	deployment := pkg.NewDeployment(config)

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
			fmt.Println("Deployment failed")
			break
		}

		if event.Value == io.EOF {
			break
		}

		if event.Service != nil {
			fmt.Printf("%s: %v\n", event.Service.Name, event.Value)
		} else {
			fmt.Printf("global: %v\n", event.Value)
		}
	}

	return deployment.Manifest.Save()
}

// NewDeployCommand creates a new `deploy` command.
func NewDeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy the configuration",
		RunE:  runDeployCommand,
	}

	return cmd
}
