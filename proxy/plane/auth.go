package plane

import (
	"github.com/gin-gonic/gin"
	"github.com/redwebcreation/nest/context"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatus(401)
			return
		}

		id, token, err := c.MustGet("nest").(*context.Context).CloudCredentials()
		if err != nil {
			c.AbortWithStatus(500)
		}

		if username != id || password != token {
			c.AbortWithStatus(401)
			return
		}

		c.Next()
	}
}

func WithNestContext(ctx *context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("nest", ctx)

		c.Next()
	}
}
