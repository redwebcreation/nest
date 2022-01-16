package pkg

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

var (
	ErrMissingService = fmt.Errorf("missing service")
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

	// Requires is a list of services that must be running before this service.
	Requires []string `yaml:"requires"`

	// Registry to pull the image from.
	// It may be a string referencing Retrieve.Registries[%s] or a Registry.
	Registry interface{} `yaml:"registry"`
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

	for _, service := range services {
		for _, require := range service.Requires {
			if _, ok := services[require]; !ok {
				return ErrMissingService
			}
		}
	}

	*s = services

	return nil
}

type ServiceMap map[string]*Service

func (s ServiceMap) BuildDependencyPlan() [][]*Service {
	graph := NewDependencyGraph(s)
	sortedServices := make([][]*Service, 0)
	maxDepth := 0

	Walker{}.Walk(graph, func(node *Node) {
		if node.Depth > maxDepth {
			for len(sortedServices) < node.Depth {
				sortedServices = append(sortedServices, []*Service{})
			}

			maxDepth = node.Depth
		}

		sortedServices[node.Depth-1] = append(sortedServices[node.Depth-1], node.Service)
	})

	for i, j := 0, len(sortedServices)-1; i < j; i, j = i+1, j-1 {
		sortedServices[i], sortedServices[j] = sortedServices[j], sortedServices[i]
	}

	return sortedServices
}

func (s ServiceMap) hasDependent(name string) bool {
	for _, dep := range s {
		for _, require := range dep.Requires {
			if require == name {
				return true
			}
		}
	}

	return false
}
