package common

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
)

type Diagnosis struct {
	Warnings []Warning `json:"warnings"`
	Errors   []Error   `json:"errors"`
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

	rawPublicIp, _ := io.ReadAll(response.Body)

	publicIp := strings.TrimSpace(string(rawPublicIp))

	for _, service := range Config.Services {
		for _, host := range service.Hosts {
			ips, err := net.LookupIP(host)
			if err != nil {
				d.NewWarning(
					fmt.Sprintf("DNS record for %s are empty", host),
					fmt.Sprintf("Add 'A %s. %s' to your DNS records", host, publicIp),
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
					"",
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
