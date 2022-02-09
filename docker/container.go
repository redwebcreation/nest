package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/redwebcreation/nest/global"
)

func CreateContainer(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (string, error) {
	res, err := docker.ContainerCreate(ctx, config, hostConfig, networkingConfig, nil, containerName)
	if err != nil {
		return "", err
	}

	global.InternalLogger.Log(global.LevelDebug, "docker.container.create", global.NewField("name", containerName), global.NewField("id", res.ID))

	return res.ID, nil
}

func StartContainer(id string) error {
	err := docker.ContainerStart(context.Background(), id, types.ContainerStartOptions{})

	if err != nil {
		return err
	}

	global.InternalLogger.Log(global.LevelDebug, "docker.container.start", global.NewField("id", id))

	return nil
}

func GetContainerIP(id string) (string, error) {
	inspection, err := docker.ContainerInspect(context.Background(), id)
	if err != nil {
		return "", err
	}

	return inspection.NetworkSettings.IPAddress, nil
}

func RunCommand(id string, command string) error {
	ref, err := docker.ContainerExecCreate(context.Background(), id, types.ExecConfig{
		Cmd: []string{"sh", "-c", command},
	})
	if err != nil {
		return err
	}

	global.InternalLogger.Log(global.LevelDebug, "docker.container.exec", global.NewField("id", id), global.NewField("command", command))

	return docker.ContainerExecStart(context.Background(), ref.ID, types.ExecStartCheck{})
}

func ListContainers() ([]types.Container, error) {
	return docker.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("label", "cloud.usenest.deployment_id"),
		),
	})
}

func RemoveContainer(id string) error {
	err := docker.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{
		Force: true,
	})

	if err != nil {
		return err
	}

	global.InternalLogger.Log(global.LevelDebug, "docker.container.remove", global.NewField("id", id))

	return nil
}
