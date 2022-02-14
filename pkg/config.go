package pkg

// Configuration represents nest's configuration
type Configuration struct {
	Services     ServiceMap  `yaml:"services" json:"services"`
	Registries   RegistryMap `yaml:"registries" json:"registries"`
	ControlPlane struct {
		Host string `yaml:"host" json:"host"`
	} `yaml:"control_plane" json:"controlPlane"`
	Proxy struct {
		Http       string `yaml:"http" json:"http"`
		Https      string `yaml:"https" json:"https"`
		SelfSigned bool   `yaml:"self_signed" json:"selfSigned"`
	} `yaml:"proxy" json:"proxy"`
}
