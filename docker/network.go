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

	global.InternalLogger.Log(global.LevelDebug, "creating a new network", global.Fields{
		"tag":  "docker.network.create",
		"name": name,
		"id":   res.ID,
	})

	return res.ID, nil
}

func ConnectContainerToNetwork(containerID, networkID string) error {
	err := docker.NetworkConnect(context.Background(), networkID, containerID, nil)

	if err != nil {
		return err
	}

	global.InternalLogger.Log(global.LevelDebug, "connecting a container to a network", global.Fields{
		"tag":          "docker.network.connect",
		"container_id": containerID,
		"network_id":   networkID,
	})

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

	global.InternalLogger.Log(global.LevelDebug, "removing a network", global.Fields{
		"tag": "docker.network.remove",
		"id":  networkID,
	})

	return nil
}
