package command

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/redwebcreation/nest/global"
	"github.com/redwebcreation/nest/pkg"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
	"sync"
	"time"
)

func runDeployCommand(cmd *cobra.Command, args []string) error {
	config, err := pkg.Config.Resolve()
	if err != nil {
		return err
	}

	graph, err := config.Services.GroupServicesInLayers()
	if err != nil {
		return err
	}

	deployment := &pkg.Deployment{
		Id:  strconv.FormatInt(time.Now().UnixMilli(), 10),
		Bus: make(pkg.MessageBus),
		Manifest: &pkg.Manifest{
			Containers: make(map[string][]*pkg.Container),
			Networks:   make(map[string]string),
		},
	}

	for k, layer := range graph {
		fmt.Printf("Deploying layer %d/%d\n", k+1, len(graph))
		inQueue := len(layer)
		messages := make(map[string]string, inQueue)

		for _, service := range layer {
			messages[service.Name] = "idle"

			go func(service *pkg.Service) {
				err = service.Deploy(deployment, k)

				if err != nil {
					deployment.Bus <- pkg.Message{
						Service: service,
						Value:   err,
					}
				}
			}(service)
		}

		render(messages)

		for message := range deployment.Bus {
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

	err = deployment.Manifest.Save(global.ConfigHome + "/manifest.json")
	if err != nil {
		return err
	}

	fmt.Println("Cleaning up old objects")
	removed, err := cleanup(deployment)
	if err != nil {
		return err
	}

	fmt.Printf("Removed %d objects\n", removed)

	return nil
}

func cleanup(deployment *pkg.Deployment) (int, error) {
	containers, err := global.Docker.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", "cloud.usenest.deployment_id"),
		),
	})
	if err != nil {
		return 0, err
	}

	networks, err := global.Docker.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", "cloud.usenest.deployment_id"),
		),
	})

	if err != nil {
		return 0, err
	}

	var containersWg sync.WaitGroup
	var removed int

	for _, container := range containers {
		if container.Labels["cloud.usenest.deployment_id"] != deployment.Id {
			containersWg.Add(1)
			go func(container types.Container) {
				defer containersWg.Done()
				err = global.Docker.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{
					Force: true,
				})
				if err != nil {
					fmt.Println(err)
				} else {
					removed++
				}
			}(container)
		}
	}

	containersWg.Wait()

	var networksWg sync.WaitGroup

	for _, network := range networks {
		if network.Labels["cloud.usenest.deployment_id"] != deployment.Id {
			networksWg.Add(1)
			go func(network types.NetworkResource) {
				defer networksWg.Done()
				err = global.Docker.NetworkRemove(context.Background(), network.ID)
				if err != nil {
					fmt.Println(err)
				} else {
					removed++
				}
			}(network)
		}
	}

	networksWg.Wait()

	return removed, nil
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
