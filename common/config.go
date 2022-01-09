package common

import (
	"fmt"
	"strings"
)

var Config *Configuration

type Configuration struct {
	Services   ServiceMap          `yaml:"services"`
	Registries map[string]Registry `yaml:"registries"`
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
		switch service.Registry.(type) {
		case Registry:
			continue
		case string:
			if service.Registry == "" {
				continue
			}

			if _, ok := c.Registries[service.Registry.(string)]; !ok {
				return ErrRegistryNotFound
			} else {
				service.Registry = c.Registries[service.Registry.(string)]
			}
		default:
			if service.Registry == nil {
				continue
			}

			return ErrInvalidRegistry
		}
	}

	return nil
}

func (s *Service) expandFromConfig(serviceName string) {
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
