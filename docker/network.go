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

	global.LogI(global.LevelDebug, "creating a new network", global.Fields{
		"tag":  "docker.network.create",
		"name": name,
		"id":   res.ID,
	})

	return res.ID, nil
}

func ConnectContainerToNetwork(networkID, containerID string) error {
	err := Client.NetworkConnect(context.Background(), networkID, containerID, nil)

	if err != nil {
		return err
	}

	global.LogI(global.LevelDebug, "connecting a container to a network", global.Fields{
		"tag":          "docker.network.connect",
		"container_id": containerID,
		"network_id":   networkID,
	})

	return nil
}
