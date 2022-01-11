package common

import (
	"fmt"
	"strings"
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
			continue
		}

		if _, ok := service.Registry.(Registry); ok {
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

func (s *Service) ExpandFromConfig(serviceName string) {
	s.Name = serviceName

	var expandedHosts []string

	for _, host := range s.Hosts {
		// expand ~example.com into example.com and www.example.com
		if strings.HasPrefix(host, "~") {
			expandedHosts = append(expandedHosts, host[1:])
			expandedHosts = append(expandedHosts, "www."+host[1:])
		} else {
			expandedHosts = append(expandedHosts, host)
		}
	}

	s.Hosts = expandedHosts

	if s.ListeningOn == "" {
		s.ListeningOn = "80"
	} else {
		s.ListeningOn = strings.TrimPrefix(s.ListeningOn, ":")
	}
}
