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
	"sync"
)

type Event struct {
	Service *Service
	Value   any
}

type Container struct {
	ID string
	IP string
}

type Manifest struct {
	Containers map[string]*Container
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
	Config   *Configuration
	Events   chan Event
	Manifest *Manifest
}

func (d *Deployment) Run() error {
	graph, err := d.Config.Services.GroupInLayers()
	if err != nil {
		return err
	}

	for layer, services := range graph {
		d.Events <- Event{nil, fmt.Sprintf("Deploying layer %d/%d\n", layer+1, len(graph))}

		var wg sync.WaitGroup
		for _, service := range services {
			wg.Add(1)
			go func(service *Service) {
				defer wg.Done()

				pipeline := DeployPipeline{
					Deployment:      d,
					Service:         service,
					HasDependencies: layer > 0 && len(service.Requires) > 0,
				}

				if err = pipeline.Run(); err != nil {
					// todo: should stop the deployment
					d.Events <- Event{service, err}
				}
			}(service)
		}
		wg.Wait()
	}

	err = d.Manifest.Save(global.GetContainerManifestFile())
	if err != nil {
		d.Events <- Event{nil, err}

		return err
	}

	if err = d.cleanup(); err != nil {
		d.Events <- Event{nil, err}

		return err
	}

	return nil
}

type DeployPipeline struct {
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

	d.Log(io.EOF)

	return nil
}

func (d *DeployPipeline) PullImage() error {
	image := docker.Image(d.Service.Image)

	return image.Pull(func(event *docker.PullEvent) {
		d.Log(event.Status)
	}, d.Deployment.Config.Registries[d.Service.Registry])
}

func (d *DeployPipeline) CreateServiceNetwork() (string, error) {
	name := fmt.Sprintf("%s_%s", d.Service.Name, d.Deployment.Id)

	net, err := docker.CreateNetwork(name, map[string]string{
		"cloud.usenest.service":       d.Service.Name,
		"cloud.usenest.deployment_id": d.Deployment.Id,
	})

	if err != nil {
		return "", err
	}

	d.Deployment.Manifest.Networks[d.Service.Name] = net

	return net, nil
}

func (d *DeployPipeline) ConnectRequiredServices(networkId string) error {
	for _, require := range d.Service.Requires {
		err := docker.ConnectContainerToNetwork(networkId, d.Deployment.Manifest.Containers[require].ID)

		if err != nil {
			return err
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

	c, err := docker.CreateContainer(context.Background(), &container.Config{
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
	}, networking, containerName)

	if err != nil {
		return "", err
	}

	return c, nil
}

func (d *DeployPipeline) RunHooks(id string, commands []string) error {
	for _, command := range commands {
		err := docker.RunCommand(id, command)
		if err != nil {
			return err
		}

		d.Log("ran command: " + command)
	}

	return nil
}

func (d *DeployPipeline) StartContainer(containerID string) error {
	err := docker.StartContainer(containerID)
	if err != nil {
		return err
	}

	ip, err := docker.GetContainerIP(containerID)
	if err != nil {
		return err
	}

	d.Deployment.Manifest.Containers[d.Service.Name] = &Container{
		ID: containerID,
		IP: ip,
	}

	fmt.Println(d.Deployment.Manifest.Containers)

	return nil
}

// todo: refactor
func (d *Deployment) cleanup() error {
	containers, err := docker.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return err
	}

	networks, err := docker.Client.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, resource := range containers {
		if resource.Labels["cloud.usenest.deployment_id"] == d.Id {
			continue
		}

		wg.Add(1)
		go func(container types.Container) {
			defer wg.Done()
			err = docker.RemoveContainer(container.ID)

			if err != nil {
				d.Events <- Event{nil, err}
			}
		}(resource)
	}

	// Containers must be removed before networks as some containers may be attached to said networks
	wg.Wait()
	wg = sync.WaitGroup{}

	for _, resource := range networks {
		if resource.Labels["cloud.usenest.deployment_id"] == d.Id {
			continue
		}

		go func(network types.NetworkResource) {
			err = docker.RemoveNetwork(network.ID)
			if err != nil {
				d.Events <- Event{nil, err}
			}
		}(resource)
	}

	wg.Wait()

	return nil
}
