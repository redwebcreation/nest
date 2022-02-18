package plane

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redwebcreation/nest/build"
	"github.com/redwebcreation/nest/context"
	"github.com/redwebcreation/nest/deploy"
)

func New(ctx *context.Context) *gin.Engine {
	// config is already resolved at this point
	config, _ := ctx.Config()
	servicesConfig, _ := ctx.ServicesConfig()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/plane/v1/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"software": "nest",
			"version":  build.Version,
			"build":    build.Commit,
		})
	})

	router.GET("/plane/v1/server", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"commit":     config.Commit,
			"branch":     config.Branch,
			"repository": config.Repository,
			"remote":     config.RemoteURL(),
			"provider":   config.Provider,
			"config":     servicesConfig,
		})
	})

	router.GET("/plane/v1/deploy", func(context *gin.Context) {
		deployment := deploy.NewDeployment(servicesConfig, ctx.Logger(), ctx.ManifestManager(), ctx.SubnetRegistryPath())

		go func() {
			err := deployment.Start()
			if err != nil {
				deployment.Events <- deploy.Event{
					Service: nil,
					Value:   deploy.ErrDeploymentFailed,
				}
			}
		}()

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type")
		context.Header("Content-Type", "text/sseEvent-stream")
		context.Header("Cache-Control", "no-cache")
		context.Header("Connection", "keep-alive")

		var err error

		for e := range deployment.Events {
			if _, ok := e.Value.(error); ok {
				err = e.Value.(error)
				break
			}

			service := "global"

			if e.Service != nil {
				service = e.Service.Name
			}

			data, _ := json.Marshal(sseEvent{
				Kind:    "log",
				Service: service,
				Data:    fmt.Sprintf("%v", e.Value),
			})

			fmt.Fprintf(context.Writer, "data: %s\n\n", data)

			context.Writer.Flush()
		}

		context.Writer.Flush()

		if err != nil {
			fmt.Fprintf(context.Writer, "data: %s\n\n", sseEvent{
				Kind: "error",
				Data: err.Error(),
			})
			context.Writer.Flush()
		}

		if err = ctx.ManifestManager().Save(deployment.Manifest); err != nil {
			fmt.Fprintf(context.Writer, "data: %s\n\n", sseEvent{
				Kind: "error",
				Data: fmt.Sprintf("%v", err),
			})
			context.Writer.Flush()
		} else {
			fmt.Fprintf(context.Writer, "data: %s\n\n", sseEvent{
				Kind: "manifest",
				Data: deployment.Manifest,
			})
			context.Writer.Flush()
		}
	})

	return router
}

type sseEvent struct {
	Kind string `json:"kind"`

	Service string `json:"service"`

	Data interface{} `json:"data"`
}

func (e sseEvent) String() string {
	data, _ := json.Marshal(e)

	return string(data)
}
