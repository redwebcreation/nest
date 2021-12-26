package common

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/me/nest/docker"
	"github.com/me/nest/global"
)

type Message struct {
	Service *Service
	Value   interface{}
}

type Deployment struct {
	Service      *Service
	ImageVersion string
}

func (d Deployment) ContainerName() string {
	return "nest_" + d.Service.Name + "_" + d.ImageVersion
}

func (d *Deployment) Start(out chan Message) {
	image := docker.Image(d.Service.Image + ":" + d.ImageVersion)

	err := image.Pull(types.ImagePullOptions{}, func(event *docker.PullEvent) {
		out <- Message{
			Service: d.Service,
			Value:   event.Status,
		}
	})

	if err != nil {
		out <- Message{
			Service: d.Service,
			Value:   err,
		}
		return
	}

	createdAt := strconv.FormatInt(time.Now().UnixMilli(), 10)
	ref, err := global.Docker.ContainerCreate(context.Background(), &container.Config{
		Image: image.String(),
		Labels: map[string]string{
			"nest:container":     "true",
			"nest:service":       d.Service.Name,
			"nest:listening_on":  d.Service.ListeningOn,
			"nest:hosts":         strings.Join(d.Service.Hosts, ","),
			"nest:image_version": d.ImageVersion,
		},
		Env: ConvertEnv(d.Service.Env),
	}, &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}, nil, nil, d.ContainerName()+"_"+createdAt)

	if err != nil {
		out <- Message{
			Service: d.Service,
			Value:   err,
		}

		return
	}

	err = global.Docker.ContainerStart(context.Background(), ref.ID, types.ContainerStartOptions{})
	if err != nil {
		out <- Message{
			Service: d.Service,
			Value:   err,
		}
		return
	}

	for _, command := range d.Service.Prestart {
		id, err := global.Docker.ContainerExecCreate(context.Background(), ref.ID, types.ExecConfig{
			Cmd: []string{"sh", "-c", command},
		})
		if err != nil {
			out <- Message{
				Service: d.Service,
				Value:   err,
			}
			return
		}

		err = global.Docker.ContainerExecStart(context.Background(), id.ID, types.ExecStartCheck{})
		if err != nil {
			out <- Message{
				Service: d.Service,
				Value:   err,
			}
			return
		}

		out <- Message{
			Service: d.Service,
			Value:   "ran command: " + command,
		}
	}

	out <- Message{
		Service: d.Service,
		Value:   fmt.Sprintf("\033[38;2;15;210;15mdeployed\033[0m (%s)", ref.ID[0:12]),
	}

	out <- Message{
		Service: d.Service,
		Value:   io.EOF,
	}

	return
}

func StopOldContainers() (int, error) {
	containers, err := docker.GetNestContainers()
	if err != nil {
		return 0, err
	}

	var wg sync.WaitGroup

	count := 0
	lastByService := make(map[string]types.Container)

	for _, container := range containers {
		if container.Created > lastByService[container.Labels["nest:service"]].Created {
			lastByService[container.Labels["nest:service"]] = container
		}
	}

	for _, container := range containers {
		isDead := Config.Services[container.Labels["nest:service"]] == nil
		if !isDead && container.ID == lastByService[container.Labels["nest:service"]].ID {
			continue
		}

		wg.Add(1)
		go func(container types.Container) {
			defer wg.Done()
			count++

			global.Docker.ContainerStop(context.Background(), container.ID, nil)
			global.Docker.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})
		}(container)
	}

	wg.Wait()
	return count, nil
}

func ConvertEnv(env map[string]string) []string {
	var dockerEnv []string

	for k, v := range env {
		dockerEnv = append(dockerEnv, fmt.Sprintf("%s=%s", k, v))
	}

	return dockerEnv
}
