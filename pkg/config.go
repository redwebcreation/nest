package pkg

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/redwebcreation/nest/docker"
	"github.com/redwebcreation/nest/global"
)

// Configuration represents nest's configuration
type Configuration struct {
	Services   ServiceMap  `yaml:"services"`
	Registries RegistryMap `yaml:"registries"`
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

// LoadConfigFromCommit loads the configuration globally from the given commit
func LoadConfigFromCommit(commit string) error {
	reader := Locator{
		Commit: commit,
	}

	contents, err := os.ReadFile(global.LocatorConfigFile)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(contents, &reader); err != nil && err.Error() == "unknown error: remote: " {
		return ErrRepositoryNotFound
	}

	Config = &reader

	return nil
}

// LoadConfig loads the configuration globally
func LoadConfig() error {
	return LoadConfigFromCommit("")
}
