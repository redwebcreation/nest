package common

import (
	"fmt"
	"regexp"
)

type Recommendation struct {
	Title   string
	Details string
}

func (r Recommendation) String() string {
	return r.Title + ": " + r.Details
}

type Error struct {
	Title string
	Error error
}

func (e Error) String() string {
	return e.Title + ": " + e.Error.Error()
}

type Diagnosis struct {
	Recommendations []Recommendation
	Errors          []Error
	Checks          []func(*Diagnosis)
}

func DiagnoseConfiguration() *Diagnosis {
	diagnosis := &Diagnosis{
		Checks: []func(*Diagnosis){
			ValidateServicesConfiguration,
		},
	}

	for _, check := range diagnosis.Checks {
		check(diagnosis)
	}

	return diagnosis
}

func ValidateServicesConfiguration(diagnosis *Diagnosis) {
	for _, service := range Config.Services {
		if service.Image == "" {
			diagnosis.NewError(Error{
				Title: fmt.Sprintf("Service %s has no image", service.Name),
			})
		}

		if len(service.Hosts) == 0 {
			diagnosis.NewError(Error{
				Title: fmt.Sprintf("Service %s has no hosts", service.Name),
			})
		}

		for k, _ := range service.Env {
			envKeyRegex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

			if !envKeyRegex.MatchString(k) {
				diagnosis.NewError(Error{
					Title: fmt.Sprintf("Service %s has invalid env key %s", service.Name, k),
				})
			}

		}

		for _, host := range service.Hosts {
			if len(host) == 0 {
				diagnosis.NewError(Error{
					Title: fmt.Sprintf("Service %s has an empty host", service.Name),
				})
			}

			if len(host) > 255 {
				diagnosis.NewError(Error{
					Title: fmt.Sprintf("Service %s has a host longer that 255 characters (%s...)", service.Name, host[0:12]),
				})
			}
		}
	}
}

func (diagnosis *Diagnosis) NewError(err Error) {
	diagnosis.Errors = append(diagnosis.Errors, err)
}
