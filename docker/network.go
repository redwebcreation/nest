package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/redwebcreation/nest/global"
)

func CreateNetwork(name string, labels map[string]string) (string, error) {
	res, err := docker.NetworkCreate(context.Background(), name, types.NetworkCreate{
		Labels: labels,
	})
	if err != nil {
		return "", err
	}

	global.InternalLogger.Log(global.LevelDebug, "docker.network.create", global.NewField("name", name), global.NewField("id", res.ID), global.NewField("labels", labels))

	return res.ID, nil
}

func ConnectContainerToNetwork(containerID, networkID string) error {
	err := docker.NetworkConnect(context.Background(), networkID, containerID, nil)

	if err != nil {
		return err
	}

	global.InternalLogger.Log(global.LevelDebug, "docker.network.connect", global.NewField("container", containerID), global.NewField("network", networkID))

	return nil
}

func ListNetworks() ([]types.NetworkResource, error) {
	return docker.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", "cloud.usenest.deployment_id"),
		),
	})
}

func RemoveNetwork(networkID string) error {
	err := docker.NetworkRemove(context.Background(), networkID)

	if err != nil {
		return err
	}

	global.InternalLogger.Log(global.LevelDebug, "docker.network.remove", global.NewField("network", networkID))

	return nil
}
