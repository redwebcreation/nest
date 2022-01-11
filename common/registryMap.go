package common

type RegistryMap map[string]*Registry

func (r *RegistryMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var registries map[string]*Registry
	if err := unmarshal(&registries); err != nil {
		return err
	}

	for name, registry := range registries {
		registry.Name = name
	}
	
	*r = registries

	return nil
}
