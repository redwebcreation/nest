package global

import (
	"fmt"

	"github.com/docker/docker/client"
)

// Docker is the global docker client
var Docker *client.Client

func LoadDocker() (*client.Client, error) {
	docker, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, fmt.Errorf("error loading docker client: %s", err)
	}

	return docker, nil
}

func init() {
	docker, err := LoadDocker()
	if err != nil {
		panic(err)
	}

	Docker = docker
}
