package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func CreateNetwork(name string, labels map[string]string) (string, error) {
	res, err := docker.NetworkCreate(context.Background(), name, types.NetworkCreate{
		Labels: labels,
	})
	if err != nil {
		return "", err
	}

	return res.ID, nil
}

func ConnectContainerToNetwork(containerID, networkID string) error {
	return docker.NetworkConnect(context.Background(), networkID, containerID, nil)
}

func ListNetworks() ([]types.NetworkResource, error) {
	return docker.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", "cloud.usenest.deployment_id"),
		),
	})
}

func RemoveNetwork(networkID string) error {
	return docker.NetworkRemove(context.Background(), networkID)
}
