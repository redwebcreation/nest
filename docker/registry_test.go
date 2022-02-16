package docker

import (
	"encoding/base64"
	"encoding/json"
	"gotest.tools/v3/assert"
	"testing"
)

func TestToBase64(t *testing.T) {
	registry := Registry{Username: "username", Password: "password"}

	b, err := registry.ToBase64()
	assert.NilError(t, err)

	// decode base64 into text
	payload, err := base64.StdEncoding.DecodeString(b)
	assert.NilError(t, err)

	bytes, _ := json.Marshal(map[string]string{
		"username": "username",
		"password": "password",
	})
	assert.Equal(t, string(payload), string(bytes))
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
		output := test.input.URLFor(test.imageName)
		assert.Equal(t, output, test.output)
	}
}
