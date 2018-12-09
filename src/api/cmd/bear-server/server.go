package main

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
	"github.com/moocss/go-webserver/src/api/router"
	"github.com/moocss/go-webserver/src/api/router/middleware"
	"net/http"
	"strings"
	"os"
	"path/filepath"
)

// New returns a app instance
func New() *gin.Engine {
	// Set gin mode.
	// gin.SetMode("debug")

	// Create the Gin engine.
	g := gin.New()

	// Routes
	router.Load(
		// Cores
		g,
		// Middlwares
		middleware.RequestId(),
	)
	return g
}


func server(c *cli.Context) error {

	var g errgroup.Group

	g.Go(func() error {
		// 启动服务

		return nil
	})
	g.Go(func() error {
		// 健康检查

		return nil
	})

	return g.Wait()
}

func redirect(w http.ResponseWriter, req *http.Request) {
	var serverHost string = ""
	serverHost = strings.TrimPrefix(serverHost, "http://")
	serverHost = strings.TrimPrefix(serverHost, "https://")
	req.URL.Scheme = "https"
	req.URL.Host = serverHost

	w.Header().Set("Strict-Transport-Security", "max-age=31536000")

	http.Redirect(w, req, req.URL.String(), http.StatusMovedPermanently)
}

func cacheDir() string {
	const base = "golang-autocert"
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, base)
	}
	return filepath.Join(os.Getenv("HOME"), ".cache", base)
}
