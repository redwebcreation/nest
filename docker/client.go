package docker

import (
	"fmt"

	"github.com/docker/docker/client"
)

var docker *client.Client

func newDocker() (*client.Client, error) {
	d, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, fmt.Errorf("error loading docker client: %s", err)
	}

	return d, nil
}

func init() {
	d, err := newDocker()
	if err != nil {
		panic(err)
	}

	docker = d
}
