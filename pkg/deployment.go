package pkg

import (
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
)

type Deployment struct {
	Id       string
	Config   *Configuration
	Events   chan Event
	Manifest *Manifest
}

var (
	ErrDeploymentFailed = fmt.Errorf("deployment failed")
)

func NewDeployment(config *Configuration) *Deployment {
	id := strconv.FormatInt(time.Now().UnixMilli(), 10)

	return &Deployment{
		Id:       id,
		Config:   config,
		Events:   make(chan Event),
		Manifest: NewManifest(id),
	}
}

func (d *Deployment) Start() error {
	graph, err := d.Config.Services.GroupInLayers()
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
