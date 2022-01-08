package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/me/nest/global"
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

func (i Image) Pull(options types.ImagePullOptions, handler func(event *PullEvent)) error {
	events, err := global.Docker.ImagePull(context.Background(), i.String(), options)
	if err != nil {
		if strings.Contains(err.Error(), "manifest for "+i.String()+" not found") {
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
