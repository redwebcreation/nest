package global

import "github.com/docker/docker/client"

var Docker *client.Client

func init() {
	docker, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		// TODO: print a nice message saying that docker is most likely not installed
		panic(err)
	}

	Docker = docker
}
