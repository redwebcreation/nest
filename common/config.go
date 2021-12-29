package common

import (
	"github.com/me/nest/global"
	"gopkg.in/yaml.v2"
	"strings"
)

var Config *Configuration

type Configuration struct {
	Services ServiceMap `yaml:"services" json:"services"`
}

func init() {
	if !global.IsConfigResolverConfigured {
		return
	}

	contents, err := global.ConfigResolver.Get("nest.yaml")
	if err != nil {
		panic(err)
	}

	var config Configuration

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		panic(err)
	}

	Config = &config
}

func (s *Service) expandFromConfig(serviceName string) {
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
	}

	if strings.HasPrefix(s.ListeningOn, ":") {
		s.ListeningOn = s.ListeningOn[1:]
	}
}
