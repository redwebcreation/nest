package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingService = errors.New("missing service")
)

type Hooks struct {
	// Prestart is a list of commands to run before the service starts.
	Prestart []string `yaml:"prestart" json:"prestart"`
	// Poststart is a list of commands to run after the service starts.
	Poststart []string `yaml:"poststart" json:"poststart"`
	// Preclean is a list of commands to run before the service is removed by the container collector.
	Preclean []string `yaml:"preclean" json:"preclean"`
	// Postclean is a list of commands to run after the service is removed by the container collector.
	Postclean []string `yaml:"postclean" json:"postclean"`
}

func (h *Hooks) MarshalJSON() ([]byte, error) {
	if h.Prestart == nil {
		h.Prestart = []string{}
	}

	if h.Poststart == nil {
		h.Poststart = []string{}
	}

	if h.Preclean == nil {
		h.Preclean = []string{}
	}

	if h.Postclean == nil {
		h.Postclean = []string{}
	}

	type plain Hooks
	return json.Marshal((*plain)(h))
}

// Service contains the information about a service.
type Service struct {
	// Name of the service.
	Name string `yaml:"-" json:"-"`
	// Include is the path to a file containing the service config.
	Include string `yaml:"include" json:"include"`

	// Image name without a registry.
	Image string `yaml:"image" json:"image"`

	// Hosts the service responds to.
	Hosts []string `yaml:"hosts" json:"hosts"`

	// Env variables for the service.
	Env EnvMap `yaml:"env" json:"env"`

	// ListeningOn is the port the service listens on.
	ListeningOn string `yaml:"listening_on" json:"listeningOn"`

	// Hooks are commands to run during the lifecycle of the service.
	Hooks Hooks `yaml:"hooks" json:"hooks"`

	// Requires is a list of services that must be running before this service.
	Requires []string `yaml:"requires" json:"requires"`

	// Registry to pull the image from.
	// It may be a string referencing Retrieve.Registries[%s] or a Registry.
	Registry string `yaml:"registry" json:"registry"`
}

func (s *Service) MarshalJSON() ([]byte, error) {
	if s.Hosts == nil {
		s.Hosts = []string{}
	}

	if s.Env == nil {
		s.Env = EnvMap{}
	}

	if s.Requires == nil {
		s.Requires = []string{}
	}

	type plain Service
	return json.Marshal((*plain)(s))
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

func (s ServiceMap) NewGraph() (*Node, error) {
	root := &Node{}
	unresolved := map[string]bool{}

	for serviceName := range s {
		edge, err := s.graph(root, serviceName, unresolved)
		if err != nil {
			return nil, err
		}

		root.AddEdge(edge)
	}

	return root, nil
}

func (s ServiceMap) graph(parent *Node, name string, unresolved map[string]bool) (*Node, error) {
	node := Node{
		Parent:  parent,
		Service: s[name],
		Depth:   parent.Depth + 1,
	}

	unresolved[name] = true

	for _, require := range s[name].Requires {
		if unresolved[require] {
			return nil, fmt.Errorf("circular dependency detected: %s -> %s", name, require)
		}

		edge, err := s.graph(&node, require, unresolved)
		if err != nil {
			return nil, err
		}

		node.AddEdge(edge)
	}

	unresolved[name] = false

	return &node, nil
}
