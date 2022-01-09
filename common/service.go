package common

import "strings"

type Service struct {
	// The path to a file containing the service configuration.
	Include string `yaml:"include" json:"include"`

	Name string
	// Image name without a tag or registry server.
	Image string `json:"image" yaml:"image"`
	// The lists of hosts the service responds to.
	Hosts []string `json:"hosts" yaml:"hosts"`
	// Environment variables for the service.
	Env map[string]string `json:"env" yaml:"env"`
	// ListeningOn is the port the service listens on.
	ListeningOn string `json:"listeningOn" yaml:"listening_on"`
	// Prestart is a list of commands to run before the service is deployed
	Prestart []string `json:"prestart" yaml:"prestart"`
	// The registry to pull the image from.
	Registry interface{} `json:"registry" yaml:"registry"`
	// The volumes to mount for the service.
	Volumes []struct {
		// The path to mount from.
		From string `json:"from" yaml:"from"`
		// The path to mount to.
		To string `json:"to" yaml:"to"`
	} `json:"volumes" yaml:"volumes"`
}

func (s *Service) Accepts(host string) bool {
	for _, h := range s.Hosts {
		if h == host {
			return true
		}

		accepted := strings.Split(h, ".")
		comparison := strings.Split(host, ".")

		for i := range comparison {
			if accepted[i] == "*" {
				comparison[i] = "*"
				continue
			}

			if accepted[i] != comparison[i] {
				break
			}
		}

		if strings.Join(comparison, ".") == strings.Join(accepted, ".") {
			return true
		}
	}

	return false
}
