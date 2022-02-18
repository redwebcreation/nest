package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redwebcreation/nest/build"
)

func NewRouter(ctx *Context) *gin.Engine {
	// config is already resolved at this point
	config, _ := ctx.Config()
	server, _ := ctx.ServerConfig()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/api/v1/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"software": "nest",
			"version":  build.Version,
			"build":    build.Commit,
		})
	})

	router.GET("/api/v1/server", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Commit":     config.Commit,
			"branch":     config.Branch,
			"repository": config.Repository,
			"remote":     config.RemoteURL(),
			"provider":   config.Provider,
			"server":     server,
		})
	})

	router.GET("/api/v1/deploy", func(context *gin.Context) {
		deployment := NewDeployment(server, ctx.ManifestManager())

		go func() {
			err := deployment.Start()
			if err != nil {
				deployment.Events <- Event{
					Service: nil,
					Value:   ErrDeploymentFailed,
				}
			}
		}()

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type")
		context.Header("Content-Type", "text/event-stream")
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

			data, _ := json.Marshal(event{
				Kind:    "log",
				Service: service,
				Data:    fmt.Sprintf("%v", e.Value),
			})

			fmt.Fprintf(context.Writer, "data: %s\n\n", data)

			context.Writer.Flush()
		}

		context.Writer.Flush()

		if err != nil {
			fmt.Fprintf(context.Writer, "data: %s\n\n", event{
				Kind: "error",
				Data: err.Error(),
			})
			context.Writer.Flush()
		}

		if err = ctx.ManifestManager().Save(deployment.Manifest); err != nil {
			fmt.Fprintf(context.Writer, "data: %s\n\n", event{
				Kind: "error",
				Data: fmt.Sprintf("%v", err),
			})
			context.Writer.Flush()
		} else {
			fmt.Fprintf(context.Writer, "data: %s\n\n", event{
				Kind: "manifest",
				Data: deployment.Manifest,
			})
			context.Writer.Flush()
		}
	})

	return router
}

type event struct {
	Kind string `json:"kind"`

	Service string `json:"service"`

	Data interface{} `json:"data"`
}

func (e event) String() string {
	data, _ := json.Marshal(e)

	return string(data)
}
