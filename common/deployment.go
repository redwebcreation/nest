package common

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/redwebcreation/nest/docker"
	"github.com/redwebcreation/nest/global"
	"io"
	"strconv"
	"time"
)

type MessageBus chan Message

type Message struct {
	Service *Service
	Value   interface{}
}

type DeployPipeline struct {
	MessageBus MessageBus
	Service    *Service
	Id         string
}

func (d DeployPipeline) Run() error {
	err := d.PullImage()
	if err != nil {
		return err
	}

	id, err := d.CreateContainer()
	if err != nil {
		return err
	}

	err = d.RunHooks(id, d.Service.Hooks.Prestart)
	if err != nil {
		return err
	}

	err = d.StartContainer(id)
	if err != nil {
		return err
	}

	err = d.RunHooks("poststart", d.Service.Hooks.Poststart)
	if err != nil {
		return err
	}

	d.MessageBus <- Message{
		Service: d.Service,
		Value:   io.EOF,
	}

	return nil
}

func (s *Service) Deploy(bus MessageBus) error {
	id := strconv.FormatInt(time.Now().UnixMilli(), 10)

	return DeployPipeline{
		MessageBus: bus,
		Service:    s,
		Id:         id,
	}.Run()
}

func (d DeployPipeline) PullImage() error {
	image := docker.Image(d.Service.Image)

	return image.Pull(func(event *docker.PullEvent) {
		d.MessageBus <- Message{
			Service: d.Service,
			Value:   event.Status,
		}
	}, d.Service.Registry.(*docker.Registry))
}

func (d DeployPipeline) CreateContainer() (string, error) {
	c, err := global.Docker.ContainerCreate(context.Background(), &container.Config{
		Image: d.Service.Image,
		Labels: map[string]string{
			"cloud.usenest.service":       d.Service.Name,
			"cloud.usenest.deployment_id": d.Id,
		},
		Env: d.Service.Env.ToDockerEnv(),
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}, nil, nil, "nest_"+d.Service.Name+"_"+d.Service.Image+"_"+d.Id)

	if err != nil {
		return "", err
	}

	return c.ID, nil
}

func (d DeployPipeline) RunHooks(id string, commands []string) error {
	for _, command := range commands {
		ref, err := global.Docker.ContainerExecCreate(context.Background(), id, types.ExecConfig{
			Cmd: []string{"sh", "-c", command},
		})
		if err != nil {
			return err
		}

		err = global.Docker.ContainerExecStart(context.Background(), ref.ID, types.ExecStartCheck{})
		if err != nil {
			return err
		}

		d.MessageBus <- Message{
			Service: d.Service,
			Value:   "ran command: " + command,
		}
	}

	return nil
}

func (d DeployPipeline) StartContainer(id string) error {
	return global.Docker.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
}
