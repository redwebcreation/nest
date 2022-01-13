package command

import (
	"fmt"
	"github.com/redwebcreation/nest/common"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
	"strings"
	"time"
)

func runDeployCommand(cmd *cobra.Command, args []string) error {
	// re-use the previous commit
	if len(args) == 0 && common.ConfigLocator.Commit != "" {
		err := LoadConfigFromCommit(common.ConfigLocator.Commit)
		if err != nil {
			return err
		}
	}

	if len(args) == 1 {
		err := common.ConfigLocator.Git.Pull(common.ConfigLocator.Branch)
		if err != nil {
			return err
		}

		commits, err := common.ConfigLocator.Git.Commits()
		if err != nil {
			return err
		}

		var commit string

		for _, c := range commits {
			if c == args[0] || strings.HasPrefix(c, args[0]) {
				commit = c
				break
			}
		}

		if commit == "" {
			return fmt.Errorf("commit not found")
		}

		err = LoadConfigFromCommit(commit)
		if err != nil {
			return err
		}
	}

	inQueue := len(common.Config.Services)
	var messages = make(map[string]string, inQueue)
	var messageBus = make(common.MessageBus)

	fmt.Printf("Using %s to deploy services.\n\n", util.White.Fg()+common.ConfigLocator.Commit[:8]+util.Reset)

	id := strconv.FormatInt(time.Now().UnixMilli(), 10)

	for _, service := range common.Config.Services {
		messages[service.Name] = "idle"

		go func(service *common.Service) {
			err := service.Deploy(id, messageBus)

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
		Use:   "deploy [commit]",
		Short: "deploy the configuration",
		RunE:  runDeployCommand,
		Args:  cobra.RangeArgs(0, 1),
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
