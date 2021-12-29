package docker

import (
	"context"
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
