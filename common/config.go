package common

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	// JsonConfig = "nest.json"
	YamlConfig = "nest.yaml"
	YmlConfig  = "nest.yml"
)

var Config *Configuration

type ServiceMap map[string]*Service

type Configuration struct {
	Services ServiceMap `yaml:"services" json:"services"`
}

func LoadConfig(path string) (*Configuration, error) {
	_, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Configuration
	// is extension yaml
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		err = yaml.Unmarshal(bytes, &config)
		// } else if strings.HasSuffix(path, ".json") {
		// err = json.Unmarshal(bytes, &config)
	} else {
		err = fmt.Errorf("unsupported config type (%s)", path)
	}

	if err != nil {
		return nil, err
	}

	return &config, nil

}

func init() {
	configFile := []string{
		YamlConfig,
		// JsonConfig,
		YmlConfig,
	}

	for i, arg := range os.Args {
		if arg == "--using" {
			if len(os.Args) < i+1 {
				panic("--using flag requires an argument")
			}

			configFile = []string{
				// The custom config has priority over any other.
				os.Args[i+1],
				configFile[0],
				configFile[1],
			}
		}
	}

	for _, configFile := range configFile {
		config, err := LoadConfig(configFile)
		if err != nil && fmt.Sprintf("open %s: no such file or directory", configFile) == err.Error() {
			continue
		}

		if err != nil {
			panic(err)
		}

		Config = config
		return
	}

	panic("could not find a config file")
}

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

			parsedService.expandFromConfig(name)

			services[name] = parsedService
			continue
		}

		service.expandFromConfig(name)
		services[name] = service
	}

	*s = services

	return nil
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
