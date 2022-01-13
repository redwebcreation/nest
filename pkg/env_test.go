package pkg

import (
	"sort"
	"testing"
)

func TestConvertEnv(t *testing.T) {
	env := EnvMap{
		"foo": "bar",
		"baz": "qux",
	}

	converted := env.ToDockerEnv()

	sort.Strings(converted)

	if converted[0] != "baz=qux" {
		t.Errorf("Expected converted env to be 'baz=qux', got '%s'", converted[1])
	}

	if converted[1] != "foo=bar" {
		t.Errorf("Expected converted env to be 'foo=bar', got '%s'", converted[0])
	}
}
