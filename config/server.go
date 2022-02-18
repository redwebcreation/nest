package config

import (
	"github.com/c-robinson/iplib"
	"github.com/redwebcreation/nest/docker"
	"github.com/redwebcreation/nest/service"
	"net"
	"sync"
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
	Network NetworkConfiguration `yaml:"network" json:"network"`
}

type NetworkConfiguration struct {
	Policy  string   `json:"policy"` // "smallest_subnet", "/{mask size}"
	Subnets []string `json:"subnets"`
}

func (n NetworkConfiguration) Manager(registryPath string, m *sync.Mutex) *docker.Subnetter {
	var subnets []iplib.Net4

	for _, subnet := range n.Subnets {
		// todo(medic): validate subnet
		ip, cidr, _ := net.ParseCIDR(subnet)

		mask, _ := cidr.Mask.Size()
		subnets = append(subnets, iplib.NewNet4(ip, mask))
	}

	return &docker.Subnetter{
		RegistryPath: registryPath,
		Subnets:      subnets,
		Lock:         m,
	}
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
