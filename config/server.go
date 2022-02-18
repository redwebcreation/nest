package config

import (
	"github.com/redwebcreation/nest/docker"
	"github.com/redwebcreation/nest/service"
)

// ServerConfig represents nest's config
type ServerConfig struct {
	Services     service.ServiceMap `yaml:"services" json:"services"`
	Registries   RegistryMap        `yaml:"registries" json:"registries"`
	ControlPlane struct {
		Host string `yaml:"host" json:"host"`
	} `yaml:"control_plane" json:"controlPlane"`
	Proxy struct {
		HTTP       string `yaml:"http" json:"http"`
		HTTPS      string `yaml:"https" json:"https"`
		SelfSigned bool   `yaml:"self_signed" json:"selfSigned"`
	} `yaml:"proxy" json:"proxy"`
	Network NetworkOptions `yaml:"network" json:"network"`
}

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

type NetworkOptions struct {
	Ipv6 bool `yaml:"ipv6" json:"ipv6"`

	// todo: add check if pool overlaps on other subnets
	// todo: add check if pool is in range of private ip ranges
	//Pools []docker.IpRange `yaml:"pools" json:"pools"`
}

//var (
//	ErrMissingIpv6Pool = errors.New("missing ipv6 pool")
//)

func (n *NetworkOptions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain NetworkOptions
	if err := unmarshal((*plain)(n)); err != nil {
		return err
	}

	//if n.Ipv6 {
	//	if len(n.Pools) == 0 {
	//		return ErrMissingIpv6Pool
	//	}
	//} else {
	//	n.Pools = docker.DefaultIpv4Pools
	//}

	return nil
}
