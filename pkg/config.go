package pkg

import (
	"fmt"
	"github.com/redwebcreation/nest/docker"
)

// Configuration represents nest's configuration
type Configuration struct {
	Services   ServiceMap  `yaml:"services"`
	Registries RegistryMap `yaml:"registries"`
	Vault      struct {
		Enabled bool   `yaml:"enabled"`
		Path    string `yaml:"token"`
	} `yaml:"vault"`
}

var (
	ErrRegistryNotFound = fmt.Errorf("registry not found")
	ErrInvalidRegistry  = fmt.Errorf("invalid registry")
)

// UnmarshalYAML implements yaml.Unmarshaler
func (c *Configuration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Configuration
	var p plain
	if err := unmarshal(&p); err != nil {
		return err
	}

	c.Registries = p.Registries
	c.Services = p.Services

	for _, service := range c.Services {
		if service.Registry == nil {
			service.Registry = docker.Registry{}
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
