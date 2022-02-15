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
	var tests = []struct {
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

	for _, test := range tests {
		output := test.input.UrlFor(test.imageName)

		if output != test.output {
			t.Errorf("Expected %s, got %s", test.output, output)
		}
	}
}
