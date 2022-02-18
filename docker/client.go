package docker

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/redwebcreation/nest/global"
	"log"
)

type Client struct {
	client *client.Client
	logger *log.Logger
}

func (c Client) Log(level global.Level, message string, fields global.Fields) {
	// we're calling the global logger only in one place
	// makes it eas
	//ier to change / refactor
	c.logger.Print(global.NewEvent(level, message, fields))
}

func NewClient() (*Client, error) {
	d, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, fmt.Errorf("error loading docker client: %s", err)
	}

	return &Client{
		client: d,
	}, nil
}

func newDefaultClient() (*Client, error) {
	return NewClient()
}
