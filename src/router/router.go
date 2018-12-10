package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moocss/go-webserver/src/handler"
	"github.com/moocss/go-webserver/src/pkgsion"
	"github.com/moocss/go-webserver/src/routerdleware/header"
)

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Welcome to api server.",
	})
}

func versionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"source":  "https://github.com/go-impatient/go-webserver",
		"version": version.Version.String(),
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
		svcd.GET("/health", handler.HealthCheck)
		svcd.GET("/disk", handler.DiskCheck)
		svcd.GET("/cpu", handler.CPUCheck)
		svcd.GET("/ram", handler.RAMCheck)
	}

	// User API
	u := g.Group("api/v1/user")
	{
		// u.GET("/:id", user.GetUserById)
	}

	return g
}
