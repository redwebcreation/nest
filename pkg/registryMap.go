package pkg

import "github.com/redwebcreation/nest/docker"

// RegistryMap maps registry names to their respective docker.Registry structs
type RegistryMap map[string]*docker.Registry

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (r *RegistryMap) UnmarshalYAML(unmarshal func(any) error) error {
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
