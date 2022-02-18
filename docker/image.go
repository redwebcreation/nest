package docker

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redwebcreation/nest/loggy"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
)

var (
	// ErrImageNotFound is returned when the image does not exist
	ErrImageNotFound = errors.New("image not found")
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

// ImagePull pulls an image from a registry
func (c Client) ImagePull(i Image, handler func(event *PullEvent), registry *Registry) error {
	image := i.String()
	options := types.ImagePullOptions{}

	if registry != nil {
		auth, err := registry.ToBase64()
		if err != nil {
			return err
		}

		options.RegistryAuth = auth

		image = registry.URLFor(image)
	}

	events, err := c.Client.ImagePull(context.Background(), image, options)
	if err != nil {
		if strings.Contains(err.Error(), "manifest for "+image+" not found") || strings.Contains(err.Error(), "repository does not exist") {
			return ErrImageNotFound
		}

		return err
	}

	c.Log(
		loggy.DebugLevel,
		"pulled docker image",
		loggy.Fields{
			"image":    image,
			"registry": registry != nil,
			"tag":      "docker.image.pull",
		},
	)

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
