package common

import (
	"fmt"
	"github.com/redwebcreation/nest/docker"
)

var Config *Configuration

type Configuration struct {
	Services   ServiceMap  `yaml:"services"`
	Registries RegistryMap `yaml:"registries"`
}

var (
	ErrRegistryNotFound = fmt.Errorf("registry not found")
	ErrInvalidRegistry  = fmt.Errorf("invalid registry")
)

func (c *Configuration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Configuration
	var p plain
	err := unmarshal(&p)
	if err != nil {
		return err
	}

	c.Registries = p.Registries
	c.Services = p.Services

	for _, service := range c.Services {
		if service.Registry == nil {
			service.Registry = &docker.Registry{}
			continue
		}

		if _, ok := service.Registry.(*docker.Registry); ok {
			continue
		}

		if _, ok := service.Registry.(string); !ok {
			return ErrInvalidRegistry
		}

		registry, ok := c.Registries[service.Registry.(string)]
		if !ok {
			return ErrRegistryNotFound
		}

		service.Registry = registry
	}

	return nil
}
