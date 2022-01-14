package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"

	"github.com/redwebcreation/nest/global"
)

var (
	// ErrImageNotFound is returned when the image does not exist
	ErrImageNotFound = fmt.Errorf("image not found")
)

// Image represents a docker image name
type Image string

// String returns the string representation of the image
func (i Image) String() string {
	return string(i)
}

// PullEvent represents a docker image pull event
type PullEvent struct {
	// Status is the current status of the pull
	Status string `json:"status"`
}

// Pull pulls an image from a registry
func (i Image) Pull(handler func(event *PullEvent), registry Registry) error {
	image := i.String()
	options := types.ImagePullOptions{}

	if !registry.IsZero() {
		auth, err := registry.ToBase64()
		if err != nil {
			return err
		}

		options.RegistryAuth = auth

		image = registry.UrlFor(image)
	}

	events, err := global.Docker.ImagePull(context.Background(), image, options)
	if err != nil {
		if strings.Contains(err.Error(), "manifest for "+image+" not found") {
			return ErrImageNotFound
		}

		return err
	}

	decoder := json.NewDecoder(events)

	var event *PullEvent

	for {
		if err = decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		handler(event)
	}

	return nil
}
