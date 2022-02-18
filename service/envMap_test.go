package service

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestEnvMap_ForDocker(t *testing.T) {
	env := EnvMap{"foo": "bar", "baz": "qux"}.ForDocker()

	assert.DeepEqual(t, env, []string{"baz=qux", "foo=bar"})
}
