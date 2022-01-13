package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/redwebcreation/nest/docker"
	"github.com/redwebcreation/nest/global"
	"os"
)

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

func LoadConfigFromCommit(commit string) error {
	reader := ConfigLocator{
		ConfigLocatorConfig: ConfigLocatorConfig{
			Commit: commit,
		},
	}

	contents, err := os.ReadFile(global.ConfigLocatorConfigFile)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(contents, &reader); err != nil && err.Error() == "unknown error: remote: " {
		return ErrRepositoryNotFound
	}

	Config = &reader

	return nil
}

func LoadConfig() error {
	return LoadConfigFromCommit("")
}
