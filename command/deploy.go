package command

import (
	"fmt"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
	"time"
)

func runDeployCommand(cmd *cobra.Command, args []string) error {
	config, err := pkg.Config.Resolve()
	if err != nil {
		return err
	}

	id := strconv.FormatInt(time.Now().UnixMilli(), 10)
	messageBus := make(pkg.MessageBus)
	graph := config.Services.BuildDependencyPlan()

	for k, layer := range graph {
		fmt.Printf("Deploying layer %d/%d\n", k+1, len(graph))
		inQueue := len(layer)
		messages := make(map[string]string, inQueue)

		for _, service := range layer {
			messages[service.Name] = "idle"

			go func(service *pkg.Service) {
				err = service.Deploy(id, messageBus)

				if err != nil {
					messageBus <- pkg.Message{
						Service: service,
						Value:   err,
					}
				}
			}(service)
		}

		render(messages)

		for message := range messageBus {
			if _, ok := message.Value.(error); ok {
				messages[message.Service.Name] = message.Value.(error).Error()
				inQueue--

				render(messages)

				if inQueue == 0 {
					break
				}

				continue
			}

			messages[message.Service.Name] = message.Value.(string)

			render(messages)
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

func render(messages map[string]string) {
	keys := make([]string, 0, len(messages))
	for k := range messages {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s: %s\n", k, messages[k])
	}
}
