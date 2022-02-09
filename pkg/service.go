package pkg

import (
	"fmt"
	"strings"
)

var (
	ErrMissingService = fmt.Errorf("missing service")
)

// Service contains the information about a service.
type Service struct {
	// Name of the service.
	Name string `yaml:"-"`
	// Include is the path to a file containing the service configuration.
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

// ApplyDefaults sets default values and transforms certain defined patterns of a unmarshalled service.
func (s *Service) ApplyDefaults(serviceName string) {
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

// Deploy starts a deployment pipeline for the service.
// todo: refactor out unclear layer parameter (the depth of the service in the graph)
func (s *Service) Deploy(deployment *Deployment, layer int) error {
	return DeployPipeline{
		Deployment:      deployment,
		Service:         s,
		HasDependencies: layer > 0 && len(s.Requires) > 0,
	}.Run()
}
