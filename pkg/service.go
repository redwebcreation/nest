package pkg

import (
	"gopkg.in/yaml.v3"
	"strings"
)

type Service struct {
	// Name of the service.
	Name string `yaml:"-"`
	// The path to a file containing the service configuration.
	Include string `yaml:"include"`

	// Image name without a tag or registry server.
	Image string `yaml:"image"`

	// Hosts the service responds to.
	Hosts []string `yaml:"hosts"`

	// Env variables for the service.
	Env EnvMap `yaml:"env"`

	// ListeningOn is the port the service listens on.
	ListeningOn string `yaml:"listening_on"`

	// Hooks are commands to run during the lifecycle of the service.
	Hooks struct {
		// Prestart is a list of commands to run before the service starts.
		Prestart []string `yaml:"prestart"`
		// Poststart is a list of commands to run after the service starts.
		Poststart []string `yaml:"poststart"`
		// Preclean is a list of commands to run before the service is removed by the container collector.
		Preclean []string `yaml:"preclean"`
		// Postclean is a list of commands to run after the service is removed by the container collector.
		Postclean []string `yaml:"postclean"`
	} `yaml:"hooks"`

	// Registry to pull the image from.
	// It may be a string referencing Retrieve.Registries[%s] or a Registry.
	Registry interface{} `yaml:"registry"`

	// Volumes to mount for the service.
	Volumes []struct {
		// The path to mount from.
		From string `yaml:"from"`
		// The path to mount to.
		To string `yaml:"to"`
	} `yaml:"volumes"`

	// Binds from the containers to the local filesystem.
	Binds []string `yaml:"binds"`
}

func (s *Service) Normalize(serviceName string) {
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

func (s *Service) Accepts(host string) bool {
	for _, h := range s.Hosts {
		if h == host {
			return true
		}

		accepted := strings.Split(h, ".")
		comparison := strings.Split(host, ".")

		for i := range comparison {
			if accepted[i] == "*" {
				comparison[i] = "*"
				continue
			}

			if accepted[i] != comparison[i] {
				break
			}
		}

		if strings.Join(comparison, ".") == strings.Join(accepted, ".") {
			return true
		}
	}

	return false
}

type ServiceMap map[string]*Service

func (s *ServiceMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var services map[string]*Service
	if err := unmarshal(&services); err != nil {
		return err
	}

	for name, service := range services {
		if service.Include != "" {
			bytes, err := Config.Git.Read(service.Include)
			if err != nil {
				return err
			}

			var parsedService *Service

			err = yaml.Unmarshal(bytes, &parsedService)
			if err != nil {
				return err
			}

			parsedService.Normalize(name)

			services[name] = parsedService
			continue
		}

		service.Normalize(name)
		services[name] = service
	}

	*s = services

	return nil
}
