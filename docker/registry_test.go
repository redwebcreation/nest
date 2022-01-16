package docker

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

func TestRegistry_IsZero(t *testing.T) {
	var dataset = []struct {
		input  Registry
		output bool
	}{
		{
			input:  Registry{},
			output: true,
		},
		{
			input: Registry{
				Name: "default",
			},
			output: false,
		},
		{
			input: Registry{
				Host: "registry.test",
			},
			output: false,
		},
		{
			input: Registry{
				Port: "5000",
			},
			output: false,
		},
		{
			input: Registry{
				Username: "username",
			},
			output: false,
		},
		{
			input: Registry{
				Password: "password",
			},
			output: false,
		},
	}

	for _, set := range dataset {
		output := set.input.IsZero()

		if output != set.output {
			t.Errorf("Expected %t, got %t", set.output, output)
		}
	}
}
