package docker

import (
	"os/exec"
	"testing"
)

func TestImage_String(t *testing.T) {
	image := Image("nginx:1")

	if image.String() != "nginx:1" {
		t.Errorf("Expected nginx:1, got %s", image.String())
	}
}

func TestImage_Pull(t *testing.T) {
	image := Image("alpine:latest")

	out, err := exec.Command("docker", "rmi", image.String()).CombinedOutput()
	if err != nil && string(out) != "No such image: alpine:latest" {
		t.Errorf("%s: %s", err, out)
	}

	_, err = exec.Command("docker", "inspect", image.String()).CombinedOutput()
	if err == nil {
		t.Error("Expected image to not exist")
	}

	err = image.Pull(func(_ *PullEvent) {}, Registry{})
	if err != nil {
		t.Error(err)
	}

	_, err = exec.Command("docker", "inspect", image.String()).CombinedOutput()
	if err != nil {
		t.Error("Expected image to exist")
	}
}
