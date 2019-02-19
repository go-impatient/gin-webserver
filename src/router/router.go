package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/moocss/go-webserver/src/model"
	"github.com/moocss/go-webserver/src/pkg/version"
	"github.com/moocss/go-webserver/src/router/middleware"
	"github.com/moocss/go-webserver/src/service"

	sdHandler "github.com/moocss/go-webserver/src/handler/api/sd"
	userHandler "github.com/moocss/go-webserver/src/handler/api/user"
)

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Welcome to api app.",
	})
}

func versionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"source":  "https://github.com/go-impatient/go-webserver",
		"version": version.Version.String(),
	})
}

// NotFound creates a gin middleware for handling page not found.
func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		model.SendError(c)
	}
}

// Load loads the services, middlewares, routes, handlers.
func Load(s service.Service, middlewares ...gin.HandlerFunc) *gin.Engine {
	// gin app
	g := gin.New()

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	// corsConfig.AllowOrigins = []string{"http://site.com"}
	g.Use(cors.New(corsConfig))

	// 使用中间件.
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	g.Use(middleware.Handler())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middleware.Version)
	g.Use(middlewares...)

	// 404 Handler.
	g.NoRoute(NotFound())

	g.GET("/version", versionHandler)
	g.GET("/", rootHandler)

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sdHandler.HealthCheck)
		svcd.GET("/disk", sdHandler.DiskCheck)
		svcd.GET("/cpu", sdHandler.CPUCheck)
		svcd.GET("/ram", sdHandler.RAMCheck)
	}

	// v1 group
	v1Group := g.Group("/api/v1")
	{
		// authentication
		v1Group.Use(middleware.Auth())
		user := g.Group("/user")
		{
			user.POST("", func(context *gin.Context) {})
			user.DELETE("/:id", func(context *gin.Context) {})
			user.PUT("/:id", func(context *gin.Context) {})
			user.GET("", func(context *gin.Context) {})
			user.GET("/:username", userHandler.GetUser(s))
		}
	}

	// v2 group
	v2Group := g.Group("/api/v2")
	{
		v2Group.GET("", func(context *gin.Context) {})
	}

	return g
}
