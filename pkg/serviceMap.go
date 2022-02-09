package pkg

import "gopkg.in/yaml.v2"

type ServiceMap map[string]*Service

func (s *ServiceMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var services map[string]*Service
	if err := unmarshal(&services); err != nil {
		return err
	}

	for name, service := range services {
		if service.Include != "" {
			bytes, err := Locator.Read(service.Include)
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

func (s ServiceMap) GroupServicesInLayers() ([][]*Service, error) {
	graph, err := s.NewGraph()
	if err != nil {
		return nil, err
	}

	return sortNodes(graph), nil
}
