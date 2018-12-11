package server

import (
	"crypto/tls"
	"github.com/moocss/go-webserver/src/storer"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"errors"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"github.com/moocss/go-webserver/src/router"
	"github.com/moocss/go-webserver/src/bootstrap"
	"github.com/moocss/go-webserver/src/log"
	"github.com/moocss/go-webserver/src/router/middleware"
)

// Init returns a app instance
func Init() *gin.Engine {
	// Set gin mode.
	gin.SetMode(bootstrap.Conf.Core.Mode)

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

func autoTLSServer() *http.Server {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(bootstrap.Conf.Core.AutoTLS.Host),
		Cache:      autocert.DirCache(bootstrap.Conf.Core.AutoTLS.Folder),
	}
	return &http.Server{
		Addr:      	":https",
		TLSConfig: 	&tls.Config{GetCertificate: m.GetCertificate},
		Handler:  	Init(),
	}
}

func defaultTLSServer() *http.Server {
	return &http.Server{
		Addr: 			"0.0.0.0:" + bootstrap.Conf.Core.TLS.Port,
		Handler:	  Init(),
	}
}

func defaultServer() *http.Server {
	return &http.Server{
		Addr: 			"0.0.0.0:" + bootstrap.Conf.Core.Port,
		Handler:	  Init(),
	}
}

// RunHTTPServer provide run http or https protocol.
func RunHTTPServer() (err error) {
	if !bootstrap.Conf.Core.Enabled {
		log.Debug("httpd server is disabled.")
		return nil
	}

	if bootstrap.Conf.Core.AutoTLS.Enabled {
		s := autoTLSServer()
		handleSignal(s)
		log.Infof("1. Start to listening the incoming requests on https address")
		err = s.ListenAndServeTLS("", "")
	} else if bootstrap.Conf.Core.TLS.CertPath != "" && bootstrap.Conf.Core.TLS.KeyPath != "" {
		s := defaultTLSServer()
		handleSignal(s)
		log.Infof("2. Start to listening the incoming requests on https address: %s", bootstrap.Conf.Core.TLS.Port)
		err = s.ListenAndServeTLS(bootstrap.Conf.Core.TLS.CertPath, bootstrap.Conf.Core.TLS.KeyPath)
	} else {
		s := defaultServer()
		handleSignal(s)
		log.Infof("3. Start to listening the incoming requests on http address: %s", bootstrap.Conf.Core.Port)
		err = s.ListenAndServe()
	}

	return
}

// handleSignal handles system signal for graceful shutdown.
func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		log.Infof("got signal [%s], exiting apiserver now", s)
		if err := server.Close(); nil != err {
			log.Errorf("server close failed ", err)
		}

		storer.DB.Close()

		log.Infof("apiserver exited")
		os.Exit(0)
	}()
}

// PingServer
func PingServer() (err error) {
	maxPingConf := bootstrap.Conf.Core.MaxPingCount
	for i := 0; i < maxPingConf; i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get("http://localhost:" + bootstrap.Conf.Core.Port + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	err = errors.New("Cannot connect to the router.")
	return err
}

// Redirect 重定向
func Redirect(w http.ResponseWriter, req *http.Request) {
	var serverHost string = ""
	serverHost = strings.TrimPrefix(serverHost, "http://")
	serverHost = strings.TrimPrefix(serverHost, "https://")
	req.URL.Scheme = "https"
	req.URL.Host = serverHost

	w.Header().Set("Strict-Transport-Security", "max-age=31536000")

	http.Redirect(w, req, req.URL.String(), http.StatusMovedPermanently)
}

// CacheDir 缓存
func CacheDir() string {
	const base = "golang-autocert"
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, base)
	}
	return filepath.Join(os.Getenv("HOME"), ".cache", base)
}