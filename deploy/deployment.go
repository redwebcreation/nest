package deploy

import (
	"errors"
	"fmt"
	"github.com/redwebcreation/nest/config"
	"github.com/redwebcreation/nest/docker"
	"io"
	"log"
	"strconv"
	"sync"
	"time"
)

type Deployment struct {
	ID                 string
	ServicesConfig     *config.ServicesConfig
	Events             chan Event
	Manifest           *Manifest
	Logger             *log.Logger
	SubnetRegistryPath string
}

var (
	ErrDeploymentFailed = errors.New("deployment failed")
)

func NewDeployment(servicesConfig *config.ServicesConfig, logger *log.Logger, manager *Manager, subnetRegistryPath string) *Deployment {
	id := strconv.FormatInt(time.Now().UnixMilli(), 10)

	return &Deployment{
		ID:                 id,
		ServicesConfig:     servicesConfig,
		Events:             make(chan Event),
		Manifest:           manager.NewManifest(id),
		Logger:             logger,
		SubnetRegistryPath: subnetRegistryPath,
	}
}

func (d *Deployment) Start() error {
	graph, err := d.ServicesConfig.Services.GroupInLayers()
	if err != nil {
		return err
	}

	var m sync.Mutex
	dockerClient, err := docker.NewClient(d.Logger, d.ServicesConfig.Network.Manager(d.SubnetRegistryPath, &m))
	if err != nil {
		return err
	}

	var errored bool
	for layer, services := range graph {
		d.Events <- Event{nil, fmt.Sprintf("Deploying layer %d/%d", layer+1, len(graph))}

		var wg sync.WaitGroup

		for _, s := range services {
			s := s // capture loop variable
			wg.Add(1)
			go func() {
				defer wg.Done()

				pipeline := Pipeline{
					Deployment:      d,
					Docker:          dockerClient,
					Service:         s,
					HasDependencies: layer > 0 && len(s.Requires) > 0,
				}

				if err = pipeline.Run(); err != nil {
					// todo: rollback
					d.Events <- Event{s, err}

					errored = true
				}
			}()

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
