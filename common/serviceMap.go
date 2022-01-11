package common

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ServiceMap map[string]*Service

func (s *ServiceMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var services map[string]*Service
	if err := unmarshal(&services); err != nil {
		return err
	}

	for name, service := range services {
		if service.Include != "" {
			bytes, err := os.ReadFile(service.Include)
			if err != nil {
				return err
			}

			var parsedService *Service

			err = yaml.Unmarshal(bytes, &parsedService)
			if err != nil {
				return err
			}

			parsedService.ExpandFromConfig(name)

			services[name] = parsedService
			continue
		}

		service.ExpandFromConfig(name)
		services[name] = service
	}

	*s = services

	return nil
}
