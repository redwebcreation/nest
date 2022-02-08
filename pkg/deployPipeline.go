package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/redwebcreation/nest/docker"
	"github.com/redwebcreation/nest/global"
	"io"
	"os"
	"strings"
)

type MessageBus chan Message

type Message struct {
	Service *Service
	Value   interface{}
}

type Container struct {
	ID string
	IP string
}

type Manifest struct {
	Containers map[string][]*Container
	Networks   map[string]string
}

func (m Manifest) Save(path string) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0600)
}

type Deployment struct {
	Id       string
	Bus      MessageBus
	Manifest *Manifest
}

type DeployPipeline struct {
	Deployment      *Deployment
	Service         *Service
	HasDependencies bool
}

func (d DeployPipeline) Run() error {
	if d.HasDependencies {
		net, err := d.CreateServiceNetwork()
		if err != nil {
			return err
		}

		err = d.ConnectRequiredServices(net)
		if err != nil {
			return err
		}
	}

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

	err = d.RunHooks(id, d.Service.Hooks.Poststart)
	if err != nil {
		return err
	}

	d.Deployment.Bus <- Message{
		Service: d.Service,
		Value:   io.EOF,
	}

	return nil
}

func (d *DeployPipeline) PullImage() error {
	image := docker.Image(d.Service.Image)

	return image.Pull(func(event *docker.PullEvent) {
		d.Deployment.Bus <- Message{
			Service: d.Service,
			Value:   event.Status,
		}
	}, d.Service.Registry.(docker.Registry))
}

func (d *DeployPipeline) CreateServiceNetwork() (string, error) {
	name := fmt.Sprintf("%s_%s", d.Service.Name, d.Deployment.Id)

	net, err := global.Docker.NetworkCreate(context.Background(), name, types.NetworkCreate{
		Labels: map[string]string{
			"cloud.usenest.service":       d.Service.Name,
			"cloud.usenest.deployment_id": d.Deployment.Id,
		},
	})

	if err != nil {
		return "", err
	}

	d.Deployment.Manifest.Networks[d.Service.Name] = net.ID

	return net.ID, nil
}

func (d *DeployPipeline) ConnectRequiredServices(networkId string) error {
	for _, require := range d.Service.Requires {
		for _, requiredContainer := range d.Deployment.Manifest.Containers[require] {
			err := global.Docker.NetworkConnect(context.Background(), networkId, requiredContainer.ID, nil)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *DeployPipeline) CreateContainer() (string, error) {
	containerName := "nest_" + d.Service.Name + "_" + strings.Replace(d.Service.Image, ":", "_", 1) + "_" + d.Deployment.Id

	var networking *network.NetworkingConfig

	if d.HasDependencies {
		networking = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				d.Service.Name: {
					NetworkID: d.Deployment.Manifest.Networks[d.Service.Name],
				},
			},
		}
	}

	c, err := global.Docker.ContainerCreate(context.Background(), &container.Config{
		Image: d.Service.Image,
		Labels: map[string]string{
			"cloud.usenest.service":       d.Service.Name,
			"cloud.usenest.deployment_id": d.Deployment.Id,
		},
		Env: d.Service.Env.ForDocker(),
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}, networking, nil, containerName)

	if err != nil {
		return "", err
	}

	return c.ID, nil
}

func (d *DeployPipeline) RunHooks(id string, commands []string) error {
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

		d.Deployment.Bus <- Message{
			Service: d.Service,
			Value:   "ran command: " + command,
		}
	}

	return nil
}

func (d *DeployPipeline) StartContainer(containerID string) error {
	err := global.Docker.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	inspection, err := global.Docker.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return err
	}

	net := "bridge"

	if d.HasDependencies {
		net = d.Service.Name
	}

	d.Deployment.Manifest.Containers[d.Service.Name] = append(d.Deployment.Manifest.Containers[d.Service.Name], &Container{
		ID: inspection.ID,
		IP: inspection.NetworkSettings.Networks[net].IPAddress,
	})

	return nil
}
