package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {

	username, password, ok := c.Request.BasicAuth()
	issuerName := os.Getenv("USER_NAME")
	issuerPassword := os.Getenv("PASSWORD")

	if !ok || username != issuerName || password != issuerPassword {
		c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	c.Next()
}
