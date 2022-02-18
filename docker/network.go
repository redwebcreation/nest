package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	logger2 "github.com/redwebcreation/nest/pkg/logger"
)

func (c Client) NetworkCreate(name string, labels map[string]string) (string, error) {
	res, err := c.client.NetworkCreate(context.Background(), name, types.NetworkCreate{
		CheckDuplicate: true,
		Labels:         labels,
	})
	if err != nil {
		return "", err
	}

	c.Log(logger2.DebugLevel, "creating a new network", logger2.Fields{
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

	c.Log(logger2.DebugLevel, "connecting a container to a network", logger2.Fields{
		"tag":          "docker.network.connect",
		"container_id": containerID,
		"network_id":   networkID,
	})

	return nil
}
