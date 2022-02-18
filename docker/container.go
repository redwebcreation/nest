package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	logger2 "github.com/redwebcreation/nest/pkg/logger"
)

func (c Client) ContainerCreate(config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (string, error) {
	res, err := c.client.ContainerCreate(context.Background(), config, hostConfig, networkingConfig, nil, containerName)
	if err != nil {
		return "", err
	}

	c.Log(logger2.DebugLevel, "creating a new docker container", logger2.Fields{
		"name": containerName,
		"id":   res.ID,
		"tag":  "docker.container.create",
	})

	return res.ID, nil
}

func (c Client) ContainerStart(id string) error {
	err := c.client.ContainerStart(context.Background(), id, types.ContainerStartOptions{})

	if err != nil {
		return err
	}

	c.Log(logger2.DebugLevel, "starting a new docker container", logger2.Fields{
		"id":  id,
		"tag": "docker.container.start",
	})

	return nil
}

func (c Client) GetContainerIP(id string) (string, error) {
	inspection, err := c.client.ContainerInspect(context.Background(), id)
	if err != nil {
		return "", err
	}

	return inspection.NetworkSettings.IPAddress, nil
}

func (c Client) ContainerExec(id string, command string) error {
	ref, err := c.client.ContainerExecCreate(context.Background(), id, types.ExecConfig{
		Cmd: []string{"sh", "-c", command},
	})
	if err != nil {
		return err
	}

	c.Log(logger2.DebugLevel, "executing a command in a container", logger2.Fields{
		"id":      id,
		"command": command,
		"tag":     "docker.container.exec",
	})

	return c.client.ContainerExecStart(context.Background(), ref.ID, types.ExecStartCheck{})
}
