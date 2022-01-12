package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"io"
	"strings"

	"github.com/redwebcreation/nest/global"
)

type Image string

func (i Image) String() string {
	return string(i)
}

type PullEvent struct {
	Status         string `json:"status"`
	Error          string `json:"error"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

func (i Image) Pull(handler func(event *PullEvent), registry *Registry) error {
	image := i.String()
	options := types.ImagePullOptions{}

	if registry != nil {
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
			return fmt.Errorf("image %s not found", i.String())
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
