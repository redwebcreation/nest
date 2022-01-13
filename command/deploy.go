package command

import (
	"fmt"
	"github.com/redwebcreation/nest/pkg"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
	"strings"
	"time"
)

func runDeployCommand(cmd *cobra.Command, args []string) error {
	config, err := pkg.Config.Retrieve()
	if err != nil {
		return err
	}

	// re-use the previous commit
	if len(args) == 0 && pkg.Config.Commit != "" {
		err := LoadConfigFromCommit(pkg.Config.Commit)
		if err != nil {
			return err
		}
	}

	if len(args) == 1 {
		err := pkg.Config.Git.Pull(pkg.Config.Branch)
		if err != nil {
			return err
		}

		commits, err := pkg.Config.Git.Commits()
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

	inQueue := len(config.Services)
	var messages = make(map[string]string, inQueue)
	var messageBus = make(pkg.MessageBus)

	fmt.Printf("Using %s to deploy services.\n\n", util.White.Fg()+pkg.Config.Commit[:8]+util.Reset)

	id := strconv.FormatInt(time.Now().UnixMilli(), 10)

	for _, service := range config.Services {
		messages[service.Name] = "idle"

		go func(service *pkg.Service) {
			err := service.Deploy(id, messageBus)

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
