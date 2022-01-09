package common

import (
	"fmt"
	"strings"

	"github.com/me/nest/global"
	"gopkg.in/yaml.v2"
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

func init() {
	if !global.IsConfigLocatorConfigured {
		return
	}

	contents, err := global.ConfigLocatorConfig.Read("nest.yaml")
	if err != nil {
		panic(err)
	}

	var config Configuration

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		panic(err)
	}

	Config = &config
}

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
		// if service.Registry is of type Registry
		switch service.Registry.(type) {
		case Registry:
		// do nothing
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
