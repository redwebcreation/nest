package deploy

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/redwebcreation/nest/docker"
	"github.com/redwebcreation/nest/service"
	"strings"
	"time"
)

type Event struct {
	Service *service.Service
	Value   any
}

type Pipeline struct {
	Docker          *docker.Client
	Deployment      *Deployment
	Service         *service.Service
	HasDependencies bool
}

func (d *Pipeline) Log(v any) {
	d.Deployment.Events <- Event{d.Service, v}
}

func (d Pipeline) Run() error {
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

	err = d.EnsureContainerIsRunning(id)
	if err != nil {
		if err2 := d.Docker.ContainerDelete(id); err2 != nil {
			return fmt.Errorf("%s (cleanup failed: %s)", err, err2)
		}

		return err
	}

	d.Log("deployment ended")

	return nil
}

func (d *Pipeline) PullImage() error {
	image := docker.Image(d.Service.Image)

	return d.Docker.ImagePull(image, func(event *docker.PullEvent) {
		d.Log(event.Status)
	}, d.Deployment.ServicesConfig.Registries[d.Service.Registry])
}

func (d *Pipeline) CreateServiceNetwork() (string, error) {
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

func (d *Pipeline) ConnectRequiredServices(networkID string) error {
	for _, require := range d.Service.Requires {
		err := d.Docker.NetworkConnect(networkID, d.Deployment.Manifest.Containers[require], []string{require})

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Pipeline) CreateContainer() (string, error) {
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

	return d.Docker.ContainerCreate(&container.Config{
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
}

func (d *Pipeline) RunHooks(containerID string, commands []string) error {
	for _, command := range commands {
		err := d.Docker.ContainerExec(containerID, command)
		if err != nil {
			return err
		}

		d.Log("ran command: " + command)
	}

	return nil
}

func (d *Pipeline) StartContainer(containerID string) error {
	err := d.Docker.ContainerStart(containerID)
	if err != nil {
		return err
	}

	d.Deployment.Manifest.Containers[d.Service.Name] = containerID

	return nil
}

// EnsureContainerIsRunning will wait for the container to start and then return
// an error if the container is not running after either :
// - 10 seconds if the container has no health-check
// - Retries * (Interval + Timeout) if the container has a health-check
//
// todo(pipeline): return logs from failed container
func (d *Pipeline) EnsureContainerIsRunning(containerID string) error {
	info, err := d.Docker.Client.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return err
	}

	var timeout time.Duration

	if info.Config.Healthcheck == nil {
		timeout = 10 * time.Second
	} else {
		seconds := float64(info.Config.Healthcheck.Retries) * (info.Config.Healthcheck.Interval + info.Config.Healthcheck.Timeout).Seconds()

		timeout = time.Duration(seconds) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("container %s is not running (timed out)", containerID)
		default:
			info, err = d.Docker.Client.ContainerInspect(ctx, containerID)
			if err != nil {
				return err
			}

			if info.State.Health != nil {
				if info.State.Health.Status == types.Healthy {
					return nil
				}

				if info.State.Health.Status == types.Unhealthy {
					return fmt.Errorf("container %s is not running (unhealthy)", containerID)
				}

				if info.State.Health.Status == types.Starting {
					continue
				}
			}

			if info.State.Status != "running" || info.RestartCount > 0 || info.State.ExitCode != 1 {
				return fmt.Errorf("container %s is not running", containerID)
			}

			time.Sleep(time.Millisecond * 500)
		}
	}
}
