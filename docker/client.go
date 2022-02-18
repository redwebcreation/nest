package docker

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/redwebcreation/nest/loggy"
	"log"
)

type Client struct {
	client        *client.Client
	logger        *log.Logger
	networkConfig *Subnetter
}

func (c Client) Log(level loggy.Level, message string, fields loggy.Fields) {
	c.logger.Print(loggy.NewEvent(level, message, fields))
}

func NewClient(logger *log.Logger, networkConf *Subnetter) (*Client, error) {
	d, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error loading docker client: %s", err)
	}

	// todo: ipv6
	return &Client{
		client:        d,
		logger:        logger,
		networkConfig: networkConf,
	}, nil
}
