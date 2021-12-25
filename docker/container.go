package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/me/nest/global"
)

func GetNestContainers() ([]types.Container, error) {
	containers, err := global.Docker.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "label",
				Value: "nest:container=true",
			},
		),
	})

	return containers, err
}

func GetServiceMap() (map[string][]types.Container, error) {
	containers, err := GetNestContainers()
	if err != nil {
		return nil, err
	}

	var nestContainers = make(map[string][]types.Container)

	for _, container := range containers {
		service := container.Labels["nest:service"]

		if service == "" || strings.HasPrefix(service, "@") {
			continue
		}

		nestContainers[service] = append(nestContainers[service], container)
	}

	return nestContainers, nil
}
