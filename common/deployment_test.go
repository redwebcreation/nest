package common

import (
	"sort"
	"testing"
)

func TestDeployment_ContainerName(t *testing.T) {
	d := Deployment{
		ImageVersion: "1.2.3",
		Service: &Service{
			Name: "test_service",
		},
	}

	if d.ContainerName() != "nest_test_service_1.2.3" {
		t.Errorf("Expected container name to be 'nest_test_service_1.2.3', got '%s'", d.ContainerName())
	}
}

func TestConvertEnv(t *testing.T) {
	env := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	converted := ConvertEnv(env)

	sort.Strings(converted)

	if converted[0] != "baz=qux" {
		t.Errorf("Expected converted env to be 'baz=qux', got '%s'", converted[1])
	}

	if converted[1] != "foo=bar" {
		t.Errorf("Expected converted env to be 'foo=bar', got '%s'", converted[0])
	}
}
