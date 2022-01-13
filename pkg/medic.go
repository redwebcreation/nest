package pkg

import (
	"fmt"
	"regexp"
)

type Diagnosis struct {
	Config   *Configuration `json:"-"`
	Warnings []Warning      `json:"warnings"`
	Errors   []Error        `json:"errors"`
}
type Warning struct {
	Title  string `json:"title"`
	Advice string `json:"advice,omitempty"`
}
type Error struct {
	Title string `json:"title"`
	Error error  `json:"error,omitempty"`
}

func DiagnoseConfiguration() *Diagnosis {
	config, err := Config.Retrieve()
	if err != nil {
		return &Diagnosis{
			Config: config,
			Errors: []Error{
				{
					Title: "Unable to load configuration",
					Error: err,
				},
			},
		}
	}

	diagnosis := Diagnosis{
		Config: config,
	}

	diagnosis.ValidateServicesConfiguration()

	return &diagnosis
}

func (d *Diagnosis) ValidateServicesConfiguration() {
	for _, service := range d.Config.Services {
		if service.Image == "" {
			d.Errors = append(d.Errors, Error{
				Title: fmt.Sprintf("Service %s has no image", service.Name),
			})
		}

		re := regexp.MustCompile(`^[a-zA-Z0-9-]+:([a-zA-Z0-9]+)$`)
		if !re.MatchString(service.Image) {
			d.Errors = append(d.Errors, Error{
				Title: fmt.Sprintf("Service %s has an invalid image", service.Name),
				Error: fmt.Errorf("image %s is not in the format <repository>:<tag>", service.Image),
			})
		} else {
			tag := re.FindStringSubmatch(service.Image)[1]

			if tag == "latest" {
				d.Errors = append(d.Errors, Error{
					Title: fmt.Sprintf("Service %s uses the `latest` tag, use a specific tag instead.", service.Name),
				})
			}
		}

		if len(service.Hosts) == 0 {
			d.Errors = append(d.Errors, Error{
				Title: fmt.Sprintf("Service %s has no hosts", service.Name),
			})
		}

		for k := range service.Env {
			envKeyRegex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

			if !envKeyRegex.MatchString(k) {
				d.Errors = append(d.Errors, Error{
					Title: fmt.Sprintf("Service %s has invalid env key %s", service.Name, k),
				})
			}
		}

		for _, host := range service.Hosts {
			if len(host) == 0 {
				d.Errors = append(d.Errors, Error{
					Title: fmt.Sprintf("Service %s has an empty host", service.Name),
				})
			}
		}
	}
}
