package config

import (
	"github.com/c-robinson/iplib"
	"github.com/redwebcreation/nest/docker"
	"github.com/redwebcreation/nest/service"
	"gopkg.in/yaml.v2"
	"net"
	"sync"
)

// ServicesConfig represents nest's config
type ServicesConfig struct {
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

func (c *ServicesConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain ServicesConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	if c.Proxy.HTTP == "" {
		c.Proxy.HTTP = "80"
	}

	if c.Proxy.HTTPS == "" {
		c.Proxy.HTTPS = "443"
	}

	return nil
}

func (c *ServicesConfig) ExpandIncludes(config *Config) error {
	for _, s := range c.Services {
		if s.Include == "" {
			continue
		}

		bytes, err := config.Read(s.Include)
		if err != nil {
			return err
		}

		var parsedService *service.Service

		err = yaml.Unmarshal(bytes, &parsedService)
		if err != nil {
			return err
		}

		parsedService.ApplyDefaults(s.Name)

		c.Services[s.Name] = parsedService
	}

	return nil
}

type NetworkConfiguration struct {
	// todo: implement smallest_subnet policy once subnetter is thoroughly tested
	Policy  string   `yaml:"policy" json:"policy"` // "smallest_subnet", "/{mask size}"
	Subnets []string `yaml:"subnets" json:"subnets"`
}

var defaultIpv4Pool = []string{"10.0.0.0/8"}

func (n *NetworkConfiguration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain NetworkConfiguration
	if err := unmarshal((*plain)(n)); err != nil {
		return err
	}

	if n.Policy == "" {
		n.Policy = "/24"
	}

	if n.Subnets == nil {
		n.Subnets = defaultIpv4Pool
	}

	return nil
}

func (n *NetworkConfiguration) Manager(registryPath string, m *sync.Mutex) *docker.Subnetter {
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
