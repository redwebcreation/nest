package pkg

import (
	"github.com/redwebcreation/nest/docker"
	"gopkg.in/yaml.v2"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestService_AssignsRegistriesToTheirServices(t *testing.T) {
	var config Configuration
	err := yaml.Unmarshal([]byte(strings.TrimSpace(`
services:
  example:
    registry: default
registries:
  default:
    host: localhost
    username: hello
    password: world`)), &config)
	if err != nil {
		t.Fatal(err)
	}

	serviceRegistry := config.Services["example"].Registry.(*docker.Registry)

	if serviceRegistry.Name != "default" {
		t.Errorf("Expected registry name to be 'default', got %s", serviceRegistry.Name)
	}

	if serviceRegistry.Host != "localhost" {
		t.Errorf("Expected registry host to be 'localhost', got %s", serviceRegistry.Host)
	}

	if serviceRegistry.Username != "hello" {
		t.Errorf("Expected registry username to be 'hello', got %s", serviceRegistry.Username)
	}

	if serviceRegistry.Password != "world" {
		t.Errorf("Expected registry password to be 'world', got %s", serviceRegistry.Password)
	}
}

func TestService_ExpandFromConfig(t *testing.T) {
	tests := []struct {
		serviceName string
		input       Service
		output      Service
	}{
		// test that the service name is set correctly
		{
			serviceName: "example",
			input:       Service{},
			output: Service{
				Name: "example",
			},
		},
		// test that tilde expansion works on hosts
		{
			serviceName: "-",
			input: Service{
				Hosts: []string{
					"~example.com",
				},
			},
			output: Service{
				Hosts: []string{
					"example.com",
					"www.example.com",
				},
			},
		},
		// test that the default port to forward the load to is 80
		{
			serviceName: "-",
			input:       Service{},
			output: Service{
				ListeningOn: "80",
			},
		},
		// test that leading : are trimmed
		{
			serviceName: "-",
			input: Service{
				ListeningOn: ":8080",
			},
			output: Service{
				ListeningOn: "8080",
			},
		},
	}

	for _, test := range tests {
		test.input.ApplyDefaults(test.serviceName)

		if test.serviceName != "-" && test.input.Name != test.serviceName {
			t.Errorf("Expected service name to be %s, got %s", test.serviceName, test.input.Name)
		}

		if test.output.ListeningOn != "" && test.input.ListeningOn != test.output.ListeningOn {
			t.Errorf("Expected listening on to be '%s', got '%s'", test.output.ListeningOn, test.input.ListeningOn)
		}

		sort.Strings(test.input.Hosts)
		sort.Strings(test.output.Hosts)

		if reflect.DeepEqual(test.input.Hosts, test.output.Hosts) == false {
			t.Errorf("Expected hosts to be %s, got %s", test.output.Hosts, test.input.Hosts)
		}

	}
}
