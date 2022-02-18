package pkg

import (
	"fmt"
	"github.com/redwebcreation/nest/docker"
	"github.com/redwebcreation/nest/pkg/loggy"
	"github.com/redwebcreation/nest/pkg/manifest"
	"io"
	"strconv"
	"sync"
	"time"
)

type Deployment struct {
	ID       string
	Server   *ServerConfig
	Events   chan Event
	Manifest *manifest.Manifest
}

var (
	ErrDeploymentFailed = fmt.Errorf("deployment failed")
)

func NewDeployment(server *ServerConfig, manager *manifest.Manager) *Deployment {
	id := strconv.FormatInt(time.Now().UnixMilli(), 10)

	return &Deployment{
		ID:       id,
		Server:   server,
		Events:   make(chan Event),
		Manifest: manager.NewManifest(id),
	}
}

func (d *Deployment) Start() error {
	graph, err := d.Server.Services.GroupInLayers()
	if err != nil {
		return err
	}

	//dockerClient, err := docker.NewClient(d.ServerConfig.Network.Ipv6, d.ServerConfig.Network.Pools)
	dockerClient, err := docker.NewClient(loggy.NewNullLogger())
	if err != nil {
		return err
	}

	var errored bool
	for layer, services := range graph {
		d.Events <- Event{nil, fmt.Sprintf("Deploying layer %d/%d", layer+1, len(graph))}

		var wg sync.WaitGroup

		for _, service := range services {
			wg.Add(1)
			go func(service *Service) {
				defer wg.Done()

				pipeline := DeployPipeline{
					Deployment:      d,
					Docker:          dockerClient,
					Service:         service,
					HasDependencies: layer > 0 && len(service.Requires) > 0,
				}

				if err = pipeline.Run(); err != nil {
					// todo: rollback
					d.Events <- Event{service, err}

					errored = true
				}
			}(service)

			if errored {
				break
			}
		}

		wg.Wait()

		if errored {
			break
		}
	}

	if errored {
		d.Events <- Event{nil, ErrDeploymentFailed}
	} else {
		d.Events <- Event{nil, io.EOF}
	}

	return nil
}
