package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/redwebcreation/nest/global"
)

func (c Client) NetworkCreate(name string, labels map[string]string) (string, error) {
	res, err := c.client.NetworkCreate(context.Background(), name, types.NetworkCreate{
		CheckDuplicate: true,
		Labels:         labels,
	})
	if err != nil {
		return "", err
	}

	c.Log(global.LevelDebug, "creating a new network", global.Fields{
		"tag":  "docker.network.create",
		"name": name,
		"id":   res.ID,
	})

	return res.ID, nil
}

func (c Client) NetworkConnect(networkID, containerID string, aliases []string) error {
	err := c.client.NetworkConnect(context.Background(), networkID, containerID, &network.EndpointSettings{
		Aliases: aliases,
	})
	if err != nil {
		return err
	}

	c.Log(global.LevelDebug, "connecting a container to a network", global.Fields{
		"tag":          "docker.network.connect",
		"container_id": containerID,
		"network_id":   networkID,
	})

	return nil
}
