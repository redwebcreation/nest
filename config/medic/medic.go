package medic

import (
	"fmt"
	"github.com/redwebcreation/nest/config"
	"github.com/redwebcreation/nest/service"
	"regexp"
)

// Diagnostic contains all the information about a given config diagnostic.
type Diagnostic struct {
	Config   *config.ServicesConfig `json:"-"`
	Warnings []Warning              `json:"warnings"`
	Errors   []Error                `json:"errors"`
}

type Warning struct {
	Title  string `json:"title"`
	Advice string `json:"advice,omitempty"`
}
type Error struct {
	Title string `json:"title"`
	Error error  `json:"error,omitempty"`
}

// MustPass ensures that the given Diagnostic has no errors.
func (d *Diagnostic) MustPass() error {
	if len(d.Errors) == 0 {
		return nil
	}

	return fmt.Errorf("invalid config (run `nest medic` for details)")
}

// DiagnoseConfig runs the diagnostics on the ServicesConfig.
func DiagnoseConfig(config *config.ServicesConfig) *Diagnostic {
	diagnostic := Diagnostic{
		Config: config,
	}

	for _, s := range diagnostic.Config.Services {
		diagnostic.ValidateService(s)
	}

	diagnostic.EnsureNoCircularDependenciesInServices()

	return &diagnostic
}

func (d *Diagnostic) ValidateService(service *service.Service) {
	if service.Image == "" {
		d.Errors = append(d.Errors, Error{
			Title: fmt.Sprintf("Service %s has no image", service.Name),
		})
	}

	if _, ok := d.Config.Registries[service.Registry]; !ok && service.Registry != "" {
		d.Errors = append(d.Errors, Error{
			Title: fmt.Sprintf("Service %s has an invalid registry", service.Name),
		})
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9-]+:([a-zA-Z0-9.]+)$`)
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

	// TODO: Hosts can be empty if the service is required by another service
	//if len(service.Hosts) == 0 {
	//	d.Errors = append(d.Errors, Error{
	//		Title: fmt.Sprintf("Service %s has no hosts", service.Name),
	//	})
	//}

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

		if host == d.Config.ControlPlane.Host {
			d.Errors = append(d.Errors, Error{
				Title: fmt.Sprintf("Service %s has the control plane host %s", service.Name, host),
			})
		}
	}
}

func (d *Diagnostic) EnsureNoCircularDependenciesInServices() {
	_, err := d.Config.Services.GroupInLayers()
	if err != nil {
		d.Errors = append(d.Errors, Error{
			Title: "Circular dependency detected",
			Error: err,
		})
	}
}
