package global

import (
	"fmt"

	"github.com/docker/docker/client"
)

// Docker is the global docker client
var Docker *client.Client

// NewDocker returns a docker client from the environment
func NewDocker() (*client.Client, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, fmt.Errorf("error loading docker client: %s", err)
	}

	return docker, nil
}

func init() {
	docker, err := NewDocker()
	if err != nil {
		panic(err)
	}

	Docker = docker
}
