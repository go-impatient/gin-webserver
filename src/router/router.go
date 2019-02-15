package router

import (
	"github.com/gin-contrib/cors"
	"net/http"

	"github.com/gin-gonic/gin"
	sdHandler "github.com/moocss/go-webserver/src/handler/api/sd"
	// userHandler "github.com/moocss/go-webserver/src/handler/api/user"
	"github.com/moocss/go-webserver/src/pkg/version"
	"github.com/moocss/go-webserver/src/router/middleware"
	errorModel "github.com/moocss/go-webserver/src/schema"
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

// NotFound creates a gin middleware for handling page not found.
func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, &errorModel.Error{
			Error:            http.StatusText(http.StatusNotFound),
			ErrorCode:        http.StatusNotFound,
			ErrorDescription: "page not found",
		})
	}
}

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mws ...gin.HandlerFunc) *gin.Engine {
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
	g.Use(mws...)

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
			user.GET("/:username", func(context *gin.Context) {})
		}
	}

	// v2 group
	v2Group := g.Group("/api/v2")
	{
		v2Group.GET("", func(context *gin.Context) {})
	}

	return g
}
