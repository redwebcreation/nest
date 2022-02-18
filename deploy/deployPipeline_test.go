package deploy

import (
	"context"
	"github.com/c-robinson/iplib"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/redwebcreation/nest/config"
	"github.com/redwebcreation/nest/docker"
	"github.com/redwebcreation/nest/loggy"
	"github.com/redwebcreation/nest/service"
	"gotest.tools/v3/assert"
	"net"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func newTestDockerClient(t *testing.T, subnetRegistryPath string) *docker.Client {
	client, err := docker.NewClient(loggy.NewNullLogger(), &docker.Subnetter{
		Lock: &sync.Mutex{},
		Subnets: []iplib.Net4{
			iplib.NewNet4(net.IP{10, 0, 0, 0}, 8),
		},
		RegistryPath: subnetRegistryPath,
	})
	assert.NilError(t, err)

	return client
}

func newTestPipeline(t *testing.T, opts ...func(*testing.T, *Pipeline)) *Pipeline {
	id := strconv.Itoa(time.Now().Nanosecond())
	path, err := os.MkdirTemp("", "subnetter")
	assert.NilError(t, err)

	pipeline := &Pipeline{
		Deployment: &Deployment{
			ID:     id,
			Events: make(chan Event),
			Logger: loggy.NewNullLogger(),
			Manifest: &Manifest{
				ID:         id,
				Containers: map[string]string{},
				Networks:   map[string]string{},
			},
			ServicesConfig:     &config.ServicesConfig{},
			SubnetRegistryPath: path,
		},
		Docker: newTestDockerClient(t, path),
		Service: &service.Service{
			Name:  "testing_service_",
			Image: "nginx:1.21.5",
		},
	}

	for _, opt := range opts {
		opt(t, pipeline)
	}

	t.Cleanup(func() {
		containers, err := pipeline.Docker.Client.ContainerList(context.Background(), types.ContainerListOptions{})
		assert.NilError(t, err)

		for _, c := range containers {
			if c.Names[0] == "/"+pipeline.ContainerName() {
				err = pipeline.Docker.ContainerDelete(c.ID)
				assert.NilError(t, err)
			}
		}

	})

	return pipeline
}

func TestPipeline_EnsureContainerIsRunning(t *testing.T) {
	pipeline := newTestPipeline(t)

	err := pipeline.EnsureContainerIsRunning("invalid_id")
	assert.Error(t, err, "Error: No such container: invalid_id")
}

func TestPipeline_EnsureContainerIsRunning2(t *testing.T) {
	pipeline := newTestPipeline(t)

	id, err := pipeline.CreateContainer()
	assert.NilError(t, err)

	err = pipeline.StartContainer(id)
	assert.NilError(t, err)

	err = pipeline.EnsureContainerIsRunning(id)
	assert.NilError(t, err)
}

func TestPipeline_EnsureContainerIsRunning3(t *testing.T) {
	pipeline := newTestPipeline(t)

	// we create the container but we don't start it
	id, err := pipeline.CreateContainer()
	assert.NilError(t, err)

	err = pipeline.EnsureContainerIsRunning(id)
	assert.ErrorIs(t, err, ErrContainerNotRunning)
}

func TestPipeline_Run(t *testing.T) {
	pipeline := newTestPipeline(t, func(t *testing.T, pipeline *Pipeline) {})

	go func() {
		for range pipeline.Deployment.Events {

		}
	}()

	err := pipeline.Run()
	assert.NilError(t, err)

	containers, err := pipeline.Docker.Client.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.Arg("name", pipeline.ContainerName()),
		),
	})
	assert.NilError(t, err)

	assert.Equal(t, len(containers), 1)
	assert.Equal(t, containers[0].State, "running")
	assert.Equal(t, containers[0].Image, "nginx:1.21.5")
	assert.Equal(t, containers[0].Labels["cloud.usenest.deployment_id"], pipeline.Deployment.ID)
	assert.Equal(t, containers[0].Labels["cloud.usenest.service"], pipeline.Service.Name)
}
