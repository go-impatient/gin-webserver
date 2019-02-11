package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moocss/go-webserver/src/pkg/version"
	"github.com/moocss/go-webserver/src/router/middleware/header"
	"github.com/moocss/go-webserver/src/router/middleware/sd"
)

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Welcome to api server.",
	})
}

func versionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"source":  "https://github.com/go-impatient/go-webserver",
		"version": version.Info.String(),
	})
}

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, middleware ...gin.HandlerFunc) *gin.Engine {
	// 使用中间件.
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	g.Use(header.NoCache)
	g.Use(header.Options)
	g.Use(header.Secure)
	g.Use(header.Version)
	g.Use(middleware...)

	g.GET("/version", versionHandler)
	g.GET("/", rootHandler)

	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "网址被外星人劫持了~")
	})

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	// User API
	u := g.Group("api/v1/user")
	{
		u.GET("/:id", func(context *gin.Context) {})
	}

	return g
}
