package command

import (
	"fmt"
	"github.com/redwebcreation/nest/common"
	"github.com/spf13/cobra"
	"sort"
)

func runDeployCommand(cmd *cobra.Command, args []string) error {
	var queued map[string]*common.Service

	if len(args) == 0 {
		queued = common.Config.Services
	} else {
		next := common.Config.Services[args[0]]

		if next == nil {
			return common.ErrServiceNotFound
		}

		queued = map[string]*common.Service{args[0]: next}
	}

	inQueue := len(queued)
	var messages = make(map[string]string, inQueue)
	var messageBus = make(common.MessageBus)

	for _, service := range queued {
		messages[service.Name] = "idle"

		go func(service *common.Service) {
			err := service.Deploy(messageBus)

			if err != nil {
				messageBus <- common.Message{
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
