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

	deployment := &pkg.Deployment{
		Id:     strconv.FormatInt(time.Now().UnixMilli(), 10),
		Config: config,
		Events: make(chan pkg.Event),
		Manifest: &pkg.Manifest{
			Containers: make(map[string]*pkg.Container),
			Networks:   make(map[string]string),
		},
	}

	go func() {
		err = deployment.Run()
	}()

	if err != nil {
		return err
	}

	// select over deployment events
	for {
		select {
		case event := <-deployment.Events:
			if event.Value == io.EOF {
				break
			}

			fmt.Println(event)
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
