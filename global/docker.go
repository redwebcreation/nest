package global

import (
	"fmt"
	"github.com/docker/docker/client"
	"os"
)

var Docker *client.Client

func init() {
	docker, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Docker client could not be initialized: %s\n", err)
		os.Exit(1)
	}

	Docker = docker
}
