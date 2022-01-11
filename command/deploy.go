package command

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/redwebcreation/nest/common"
	"github.com/redwebcreation/nest/util"
	"github.com/spf13/cobra"
)

var deploymentStart int64

func runDeployCommand(cmd *cobra.Command, args []string) error {
	var deployables common.ServiceMap

	if len(args) == 0 {
		deployables = common.Config.Services
	} else {
		deployable := common.Config.Services[args[0]]

		if deployable == nil {
			return fmt.Errorf("service %s not found", args[0])
		}

		deployables = common.ServiceMap{
			deployable.Name: deployable,
		}
	}

	deploymentStart = time.Now().UnixMilli()
	deployableSize := len(deployables)
	inQueue := deployableSize
	messages := make(map[string]string, inQueue)
	updates := make(chan common.Message)

	for _, service := range deployables {
		messages[service.Name] = "idle"
		deployment := common.Deployment{
			Service:      service,
			ImageVersion: "latest",
		}

		go deployment.Start(updates)
	}

	render(messages)

	for update := range updates {
		if update.Value == io.EOF {
			inQueue--

			if inQueue == 0 {
				break
			}

			continue
		}

		if _, ok := update.Value.(error); ok {
			messages[update.Service.Name] = update.Value.(error).Error()
			inQueue--

			render(messages)

			if inQueue == 0 {
				break
			}

			continue
		} else {
			messages[update.Service.Name] = update.Value.(string)
		}

		render(messages)
	}

	fmt.Printf("\nDeployed %d %s in %.3fs.\n", deployableSize, util.Plural(deployableSize, "service", "services"), float64(time.Now().UnixMilli()-deploymentStart)/1000.0)

	stopped, _ := common.StopOldContainers()
	fmt.Printf("Cleaned up %d older services.\n", stopped)

	return nil
}

func DeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy the configuration",
		RunE:  runDeployCommand,
	}

	return WithConfig(cmd)
}

func render(updates map[string]string) {
	keys := make([]string, 0, len(updates))
	for k := range updates {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	fmt.Println("\033[H\033[2J")

	for _, k := range keys {
		fmt.Printf("%s: %s\n", k, updates[k])
	}
}