package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/redwebcreation/nest/loggy"
)

func (c Client) ContainerCreate(config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (string, error) {
	res, err := c.Client.ContainerCreate(context.Background(), config, hostConfig, networkingConfig, nil, containerName)
	if err != nil {
		return "", err
	}

	c.Log(loggy.DebugLevel, "creating a new docker container", loggy.Fields{
		"name": containerName,
		"id":   res.ID,
		"tag":  "docker.container.create",
	})

	return res.ID, nil
}

func (c Client) ContainerStart(id string) error {
	err := c.Client.ContainerStart(context.Background(), id, types.ContainerStartOptions{})

	if err != nil {
		return err
	}

	c.Log(loggy.DebugLevel, "starting a new docker container", loggy.Fields{
		"id":  id,
		"tag": "docker.container.start",
	})

	return nil
}

func (c Client) GetContainerIP(id string) (string, error) {
	info, err := c.Client.ContainerInspect(context.Background(), id)
	if err != nil {
		return "", err
	}

	return info.NetworkSettings.IPAddress, nil
}

func (c Client) ContainerExec(id string, command string) error {
	ref, err := c.Client.ContainerExecCreate(context.Background(), id, types.ExecConfig{
		Cmd: []string{"sh", "-c", command},
	})
	if err != nil {
		return err
	}

	c.Log(loggy.DebugLevel, "executing a command in a container", loggy.Fields{
		"id":      id,
		"command": command,
		"tag":     "docker.container.exec",
	})

	return c.Client.ContainerExecStart(context.Background(), ref.ID, types.ExecStartCheck{})
}

// ContainerDelete stops and removes a container
func (c Client) ContainerDelete(containerID string) error {
	err := c.Client.ContainerStop(context.Background(), containerID, nil)
	if err != nil {
		return err
	}

	c.Log(loggy.DebugLevel, "stopping a container", loggy.Fields{
		"id":  containerID,
		"tag": "docker.container.stop",
	})

	err = c.Client.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{})

	c.Log(loggy.DebugLevel, "removing a container", loggy.Fields{
		"id":  containerID,
		"tag": "docker.container.remove",
	})

	return err
}
