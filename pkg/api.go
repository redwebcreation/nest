package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/redwebcreation/nest/global"
)

func NewRouter(config *Configuration) *gin.Engine {
	router := gin.New()
	if global.Version == "dev" {
		gin.SetMode(gin.DebugMode)
		router.Use(gin.Logger())
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(gin.Recovery())

	router.GET("/api/v1/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"software": "nest",
			"version":  global.Version,
		})
	})

	router.GET("/api/v1/config", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"commit":     Locator.Commit,
			"branch":     Locator.Branch,
			"repository": Locator.Repository,
			"remote":     Locator.RemoteURL(),
			"provider":   Locator.Provider,
			"config":     config,
		})
	})

	router.GET("/api/v1/deploy", func(context *gin.Context) {

	})

	return router
}
