package command

import (
	"fmt"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"io"
	"strconv"
	"time"
)

func runDeployCommand(cmd *cobra.Command, args []string) error {
	config, err := pkg.Locator.Resolve()
	if err != nil {
		return err
	}

	id := strconv.FormatInt(time.Now().UnixMilli(), 10)
	deployment := &pkg.Deployment{
		Id:       id,
		Config:   config,
		Events:   make(chan pkg.Event),
		Manifest: pkg.NewManifest(id),
	}

	go func() {
		err = deployment.Run()
	}()

	if err != nil {
		return err
	}

	for event := range deployment.Events {
		if event.Value == io.EOF {
			break
		}

		if event.Service != nil {
			fmt.Printf("%s: %v\n", event.Service.Name, event.Value)
		} else {
			fmt.Printf("global: %v\n", event.Value)
		}

	}

	return nil
}

// NewDeployCommand creates and configures the services defined in the configuration
func NewDeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy the configuration",
		RunE:  runDeployCommand,
	}

	return cmd
}
