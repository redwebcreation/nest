package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

func CreateContainer(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (string, error) {
	res, err := docker.ContainerCreate(ctx, config, hostConfig, networkingConfig, nil, containerName)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}

func StartContainer(id string) error {
	return docker.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
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
	return docker.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{
		Force: true,
	})
}
