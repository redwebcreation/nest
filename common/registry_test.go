package common

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestToBase64(t *testing.T) {
	registry := Registry{
		Username: "username",
		Password: "password",
	}

	b, err := registry.ToBase64()
	if err != nil {
		t.Error(err)
	}

	// decode base64 into text
	payload, err := base64.StdEncoding.DecodeString(b)

	if err != nil {
		t.Error(err)
	}

	bytes, _ := json.Marshal(map[string]string{
		"username": "username",
		"password": "password",
	})

	if string(payload) != string(bytes) {
		t.Errorf("Expected %s, got %s", string(bytes), string(payload))
	}
}

func TestIsDefault(t *testing.T) {
	registry := Registry{
		Name: "whatever",
	}

	if registry.IsDefault() {
		t.Errorf("Expected registry [%s] to not be default", registry.Name)
	}

	registry.Name = ""

	if !registry.IsDefault() {
		t.Error("Expected registry [] to be default")
	}

	registry.Name = "default"

	if !registry.IsDefault() {
		t.Error("Expected registry [default] to be default")
	}

	registry.Name = "@"

	if !registry.IsDefault() {
		t.Error("Expected registry [@] to be default")
	}
}

func TestRegistry_UrlFor(t *testing.T) {
	var dataset = []struct {
		imageName string
		input     Registry
		output    string
	}{
		{
			imageName: "nginx",
			input: Registry{
				Host: "registry.test",
			},
			output: "registry.test/nginx",
		},
		{
			imageName: "nginx",
			input: Registry{
				Host: "registry.test",
				Port: "5000",
			},
			output: "registry.test:5000/nginx",
		},
	}

	for _, set := range dataset {
		output := set.input.UrlFor(set.imageName)

		if output != set.output {
			t.Errorf("Expected %s, got %s", set.output, output)
		}
	}
}
