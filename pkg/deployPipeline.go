package pkg

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/redwebcreation/nest/docker"
	"strings"
)

type Event struct {
	Service *Service
	Value   any
}

type DeployPipeline struct {
	Docker          *docker.Client
	Deployment      *Deployment
	Service         *Service
	HasDependencies bool
}

func (d *DeployPipeline) Log(v any) {
	d.Deployment.Events <- Event{d.Service, v}
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

	d.Log("deployment ended")

	return nil
}

func (d *DeployPipeline) PullImage() error {
	image := docker.Image(d.Service.Image)

	return d.Docker.ImagePull(image, func(event *docker.PullEvent) {
		d.Log(event.Status)
	}, d.Deployment.Config.Registries[d.Service.Registry])
}

func (d *DeployPipeline) CreateServiceNetwork() (string, error) {
	name := fmt.Sprintf("%s_%s", d.Service.Name, d.Deployment.ID)

	net, err := d.Docker.NetworkCreate(name, map[string]string{
		"cloud.usenest.service":       d.Service.Name,
		"cloud.usenest.deployment_id": d.Deployment.ID,
	})

	if err != nil {
		return "", err
	}

	d.Deployment.Manifest.Networks[d.Service.Name] = net

	return net, nil
}

func (d *DeployPipeline) ConnectRequiredServices(networkID string) error {
	for _, require := range d.Service.Requires {
		err := d.Docker.NetworkConnect(networkID, d.Deployment.Manifest.Containers[require], []string{require})

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DeployPipeline) CreateContainer() (string, error) {
	containerName := "nest_" + d.Service.Name + "_" + strings.Replace(d.Service.Image, ":", "_", 1) + "_" + d.Deployment.ID

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

	c, err := d.Docker.ContainerCreate(&container.Config{
		Image: d.Service.Image,
		Labels: map[string]string{
			"cloud.usenest.service":       d.Service.Name,
			"cloud.usenest.deployment_id": d.Deployment.ID,
		},
		Env: d.Service.Env.ForDocker(),
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}, networking, containerName)

	if err != nil {
		return "", err
	}

	return c, nil
}

func (d *DeployPipeline) RunHooks(id string, commands []string) error {
	for _, command := range commands {
		err := d.Docker.ContainerExec(id, command)
		if err != nil {
			return err
		}

		d.Log("ran command: " + command)
	}

	return nil
}

func (d *DeployPipeline) StartContainer(containerID string) error {
	err := d.Docker.ContainerStart(containerID)
	if err != nil {
		return err
	}

	d.Deployment.Manifest.Containers[d.Service.Name] = containerID

	return nil
}
