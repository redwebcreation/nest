package docker

import (
	"fmt"
	"github.com/docker/docker/client"
	logger2 "github.com/redwebcreation/nest/pkg/loggy"
	"log"
)

type Client struct {
	client *client.Client
	logger *log.Logger
}

func (c Client) Log(level logger2.Level, message string, fields logger2.Fields) {
	c.logger.Print(logger2.NewEvent(level, message, fields))
}

func NewClient(logger *log.Logger) (*Client, error) {
	d, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, fmt.Errorf("error loading docker client: %s", err)
	}

	return &Client{
		client: d,
		logger: logger,
	}, nil
}
