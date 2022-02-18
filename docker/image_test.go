package docker

import (
	"github.com/redwebcreation/nest/loggy"
	"gotest.tools/v3/assert"
	"os/exec"
	"strings"
	"testing"
)

func TestImage_String(t *testing.T) {
	assert.Equal(t, Image("nginx:1").String(), "nginx:1")
}

func TestClient_ImagePull(t *testing.T) {
	image := Image("alpine:latest")

	_ = exec.Command("docker", "rmi", image.String()).Run()

	out, err := exec.Command("docker", "inspect", image.String()).CombinedOutput()
	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(string(out), "No such object: alpine:latest"))

	client, err := NewClient(loggy.NewNullLogger(), nil)
	assert.NilError(t, err)

	err = client.ImagePull(image, func(_ *PullEvent) {}, nil)
	assert.NilError(t, err)

	_, err = exec.Command("docker", "inspect", image.String()).CombinedOutput()
	assert.NilError(t, err, "expected image to exist once pulled")
}
