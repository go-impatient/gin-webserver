package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/moocss/go-webserver/src/api/pkg/version"
)

// Version is a middleware function that appends the Bear version information
// to the HTTP response. This is intended for debugging and troubleshooting.
func Version(c *gin.Context) {
	c.Header("X-BEAR-VERSION", version.Version.String())
}
