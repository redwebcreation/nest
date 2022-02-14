package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/redwebcreation/nest/global"
)

func CreateNetwork(name string, labels map[string]string) (string, error) {
	res, err := Client.NetworkCreate(context.Background(), name, types.NetworkCreate{
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
	err := Client.NetworkConnect(context.Background(), networkID, containerID, nil)

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

func RemoveNetwork(networkID string) error {
	err := Client.NetworkRemove(context.Background(), networkID)

	if err != nil {
		return err
	}

	global.InternalLogger.Log(global.LevelDebug, "removing a network", global.Fields{
		"tag": "docker.network.remove",
		"id":  networkID,
	})

	return nil
}
