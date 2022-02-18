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
	Name string `yaml:"-" json:"-"`
	// Include is the path to a file containing the service config.
	Include string `yaml:"include" json:"include"`

	// Image name without a tag or registry serverConfig.
	Image string `yaml:"image" json:"image"`

	// Hosts the service responds to.
	Hosts []string `yaml:"hosts" json:"hosts"`

	// Env variables for the service.
	Env EnvMap `yaml:"env" json:"env"`

	// ListeningOn is the port the service listens on.
	ListeningOn string `yaml:"listening_on" json:"listeningOn"`

	// Hooks are commands to run during the lifecycle of the service.
	Hooks struct {
		// Prestart is a list of commands to run before the service starts.
		Prestart []string `yaml:"prestart" json:"prestart"`
		// Poststart is a list of commands to run after the service starts.
		Poststart []string `yaml:"poststart" json:"poststart"`
		// Preclean is a list of commands to run before the service is removed by the container collector.
		Preclean []string `yaml:"preclean" json:"preclean"`
		// Postclean is a list of commands to run after the service is removed by the container collector.
		Postclean []string `yaml:"postclean" json:"postclean"`
	} `yaml:"hooks" json:"hooks"`

	// Requires is a list of services that must be running before this service.
	Requires []string `yaml:"requires" json:"requires"`

	// Registry to pull the image from.
	// It may be a string referencing Retrieve.Registries[%s] or a Registry.
	Registry string `yaml:"registry" json:"registry"`
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

type ServiceMap map[string]*Service

func (s *ServiceMap) UnmarshalYAML(unmarshal func(any) error) error {
	var services map[string]*Service
	if err := unmarshal(&services); err != nil {
		return err
	}

	for name, service := range services {
		if service.Include != "" {
			// todo: handle include
			delete(services, name)
			continue
			//bytes, err := Config.Read(service.Include)
			//if err != nil {
			//	return err
			//}
			//
			//var parsedService *Service
			//
			//err = yaml.Unmarshal(bytes, &parsedService)
			//if err != nil {
			//	return err
			//}
			//
			//parsedService.ApplyDefaults(name)
			//
			//services[name] = parsedService
		}

		service.ApplyDefaults(name)
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

func (s ServiceMap) GroupInLayers() ([][]*Service, error) {
	graph, err := s.NewGraph()
	if err != nil {
		return nil, err
	}

	return sortNodes(graph), nil
}
