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
}

func DiagnoseConfiguration() *Diagnosis {
	var diagnosis Diagnosis

	diagnosis.ValidateServicesConfiguration()
	diagnosis.EnsureDnsRecordPointsToHost()

	return &diagnosis
}

func (d *Diagnosis) ValidateServicesConfiguration() {
	for _, service := range Config.Services {
		if service.Image == "" {
			d.NewError(
				fmt.Sprintf("Service %s has no image", service.Name),
				nil,
			)
		}

		if len(service.Hosts) == 0 {
			d.NewError(
				fmt.Sprintf("Service %s has no hosts", service.Name),
				nil,
			)
		}

		for k := range service.Env {
			envKeyRegex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

			if !envKeyRegex.MatchString(k) {
				d.NewError(
					fmt.Sprintf("Service %s has invalid env key %s", service.Name, k),
					nil,
				)
			}
		}

		for _, host := range service.Hosts {
			if len(host) == 0 {
				d.NewError(
					fmt.Sprintf("Service %s has an empty host", service.Name),
					nil,
				)
			}
		}
	}
}

func (d *Diagnosis) NewError(title string, err error) {
	d.Errors = append(d.Errors, Error{
		Title: title,
		Error: err,
	})
}

func (d *Diagnosis) EnsureDnsRecordPointsToHost() {
	response, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		d.NewError("Couldn't retrieve your public IP address", err)
		return
	}
	defer response.Body.Close()

	rawPublicIp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		d.NewError("Couldn't retrieve your public IP address", err)
		return
	}

	publicIp := strings.TrimSpace(string(rawPublicIp))

	for _, service := range Config.Services {
		for _, host := range service.Hosts {
			ips, err := net.LookupIP(host)
			if err != nil {
				d.NewWarning(
					fmt.Sprintf("DNS record for %s are empty", host),
					fmt.Sprintf("Create a DNS record for %s pointing to %s", host, publicIp),
				)
				continue
			}

			hasMatchingIp := false

			for _, ip := range ips {
				if ip.String() == publicIp {
					hasMatchingIp = true
				}
			}

			if !hasMatchingIp {
				d.NewWarning(
					fmt.Sprintf("DNS record for %s does not point to %s", host, publicIp),
					fmt.Sprintf("Create a DNS record for %s pointing to %s", host, publicIp),
				)
			}
		}
	}
}

func (d *Diagnosis) NewWarning(title string, advice string) {
	d.Warnings = append(d.Warnings, Warning{
		Title:  title,
		Advice: advice,
	})
}
