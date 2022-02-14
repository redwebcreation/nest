package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/redwebcreation/nest/global"
)

func CreateContainer(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (string, error) {
	res, err := Client.ContainerCreate(ctx, config, hostConfig, networkingConfig, nil, containerName)
	if err != nil {
		return "", err
	}

	global.InternalLogger.Log(global.LevelDebug, "creating a new docker container", global.Fields{
		"name": containerName,
		"id":   res.ID,
		"tag":  "docker.container.create",
	})

	return res.ID, nil
}

func StartContainer(id string) error {
	err := Client.ContainerStart(context.Background(), id, types.ContainerStartOptions{})

	if err != nil {
		return err
	}

	global.InternalLogger.Log(global.LevelDebug, "starting a new docker container", global.Fields{
		"id":  id,
		"tag": "docker.container.start",
	})

	return nil
}

func GetContainerIP(id string) (string, error) {
	inspection, err := Client.ContainerInspect(context.Background(), id)
	if err != nil {
		return "", err
	}

	return inspection.NetworkSettings.IPAddress, nil
}

func RunCommand(id string, command string) error {
	ref, err := Client.ContainerExecCreate(context.Background(), id, types.ExecConfig{
		Cmd: []string{"sh", "-c", command},
	})
	if err != nil {
		return err
	}

	global.InternalLogger.Log(global.LevelDebug, "executing a command in a container", global.Fields{
		"id":      id,
		"command": command,
		"tag":     "docker.container.exec",
	})

	return Client.ContainerExecStart(context.Background(), ref.ID, types.ExecStartCheck{})
}

func RemoveContainer(id string) error {
	err := Client.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{
		Force: true,
	})

	if err != nil {
		return err
	}

	global.InternalLogger.Log(global.LevelDebug, "removing a docker container", global.Fields{
		"id":  id,
		"tag": "docker.container.remove",
	})

	return nil
}
