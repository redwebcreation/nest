package common

import "github.com/redwebcreation/nest/docker"

type RegistryMap map[string]*docker.Registry

func (r *RegistryMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var registries map[string]*docker.Registry
	if err := unmarshal(&registries); err != nil {
		return err
	}

	for name, registry := range registries {
		registry.Name = name
	}

	*r = registries

	return nil
}
