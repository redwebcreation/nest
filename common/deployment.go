package common

import (
	"io"
)

type MessageBus chan Message

type Message struct {
	Service *Service
	Value   interface{}
}

type DeployPipeline struct {
	MessageBus MessageBus
	Service    *Service
}

func (d DeployPipeline) Run() error {
	err := d.PullImage()
	if err != nil {
		return err
	}

	id, err := d.CreateContainer()
	if err != nil {
		return err
	}

	err = d.RunHooks("prestart", id)
	if err != nil {
		return err
	}

	err = d.StartContainer(id)
	if err != nil {
		return err
	}

	err = d.RunHooks("poststart", id)
	if err != nil {
		return err
	}

	d.MessageBus <- Message{
		Service: d.Service,
		Value:   io.EOF,
	}

	return nil
}

func (s *Service) Deploy(bus MessageBus) error {
	return DeployPipeline{
		MessageBus: bus,
		Service:    s,
	}.Run()
}

func (d DeployPipeline) PullImage() error {
	return nil
}

func (d DeployPipeline) CreateContainer() (string, error) {
	return "", nil
}

func (d DeployPipeline) RunHooks(hook string, _ ...interface{}) error {
	return nil
}

func (d DeployPipeline) StartContainer(id string) error {
	return nil
}

//func (d Deployment) ContainerName() string {
//	return "nest_" + d.Service.Name + "_" + d.ImageVersion
//}
//func (d *Deployment) Start(out chan Message) {
//	image := docker.Image(d.Service.Image + ":" + d.ImageVersion)
//
//	err := image.Pull(types.ImagePullOptions{}, func(event *docker.PullEvent) {
//		out <- Message{
//			Service: d.Service,
//			Value:   event.Status,
//		}
//	})
//
//	if err != nil {
//		out <- Message{
//			Service: d.Service,
//			Value:   err,
//		}
//		return
//	}
//
//	createdAt := strconv.FormatInt(time.Now().UnixMilli(), 10)
//	ref, err := global.Docker.ContainerCreate(context.Background(), &container.Config{
//		Image: image.String(),
//		Labels: map[string]string{
//			"nest:container":     "true",
//			"nest:service":       d.Service.Name,
//			"nest:listening_on":  d.Service.ListeningOn,
//			"nest:hosts":         strings.Join(d.Service.Hosts, ","),
//			"nest:image_version": d.ImageVersion,
//		},
//		Env: ConvertEnv(d.Service.Env),
//	}, &container.HostConfig{
//		RestartPolicy: container.RestartPolicy{
//			Name: "always",
//		},
//	}, nil, nil, d.ContainerName()+"_"+createdAt)
//
//	if err != nil {
//		out <- Message{
//			Service: d.Service,
//			Value:   err,
//		}
//
//		return
//	}
//
//	err = global.Docker.ContainerStart(context.Background(), ref.ID, types.ContainerStartOptions{})
//	if err != nil {
//		out <- Message{
//			Service: d.Service,
//			Value:   err,
//		}
//		return
//	}
//
//	for _, command := range d.Service.Prestart {
//		id, err := global.Docker.ContainerExecCreate(context.Background(), ref.ID, types.ExecConfig{
//			Cmd: []string{"sh", "-c", command},
//		})
//		if err != nil {
//			out <- Message{
//				Service: d.Service,
//				Value:   err,
//			}
//			return
//		}
//
//		err = global.Docker.ContainerExecStart(context.Background(), id.ID, types.ExecStartCheck{})
//		if err != nil {
//			out <- Message{
//				Service: d.Service,
//				Value:   err,
//			}
//			return
//		}
//
//		out <- Message{
//			Service: d.Service,
//			Value:   "ran command: " + command,
//		}
//	}
//
//	out <- Message{
//		Service: d.Service,
//		Value:   fmt.Sprintf("\033[38;2;15;210;15mdeployed\033[0m (%s)", ref.ID[0:12]),
//	}
//
//	out <- Message{
//		Service: d.Service,
//		Value:   io.EOF,
//	}
//}
//
//func ConvertEnv(env map[string]string) []string {
//	var dockerEnv []string
//
//	for k, v := range env {
//		dockerEnv = append(dockerEnv, fmt.Sprintf("%s=%s", k, v))
//	}
//
//	return dockerEnv
//}
