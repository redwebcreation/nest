package common

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
)

type Warning struct {
	Title  string
	Advice string
}

func (r Warning) String() string {
	return r.Title + ": " + r.Advice
}

type Error struct {
	Title string
	Error error
}

func (e Error) String() string {
	return e.Title + ": " + e.Error.Error()
}

type Diagnosis struct {
	Warnings []Warning
	Errors   []Error
	Checks   []func(*Diagnosis)
}

func DiagnoseConfiguration() *Diagnosis {
	diagnosis := &Diagnosis{
		Checks: []func(*Diagnosis){
			ValidateServicesConfiguration,
			EnsureDnsRecordPointsToHost,
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

		for k := range service.Env {
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

func EnsureDnsRecordPointsToHost(diagnosis *Diagnosis) {
	response, err := http.Get("http://checkip.amazonaws.com")
	if err == nil {
		defer response.Body.Close()
	} else if strings.Contains(err.Error(), "no such host") {
		diagnosis.NewWarning(Warning{
			Title:  "It looks like you're not connected to internet, or AWS is down.",
			Advice: "If this is a production server, GL HF, you'll need it.",
		})
		return
	} else {
		diagnosis.NewError(Error{
			Title: "Couldn't retrieve your public IP address",
			Error: err,
		})
		return
	}

	rawPublicIp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		diagnosis.NewError(Error{
			Title: "Couldn't retrieve your public IP address",
			Error: err,
		})
		return
	}

	publicIp := strings.TrimSpace(string(rawPublicIp))

	for _, service := range Config.Services {
		for _, host := range service.Hosts {
			ips, err := net.LookupIP(host)
			if err != nil {
				diagnosis.NewWarning(Warning{
					Title:  fmt.Sprintf("DNS record for %s does not exist", host),
					Advice: fmt.Sprintf("Create a DNS record for %s pointing to %s", host, publicIp),
				})
				continue
			}

			hasMatchingIp := false

			for _, ip := range ips {
				if ip.String() == publicIp {
					hasMatchingIp = true
				}
			}

			if !hasMatchingIp {
				diagnosis.NewWarning(Warning{
					Title:  fmt.Sprintf("DNS record for %s does not point to %s", host, publicIp),
					Advice: fmt.Sprintf("Create a DNS record for %s pointing to %s", host, publicIp),
				})
			}
		}
	}
}

func (d *Diagnosis) NewWarning(w Warning) {
	d.Warnings = append(d.Warnings, w)
}
