package main

import (
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/moocss/apiserver/src/router"
	"github.com/moocss/apiserver/src/router/middleware"
	"github.com/moocss/apiserver/src/service"
	"golang.org/x/crypto/acme/autocert"
)

// New returns a app instance
func New() *gin.Engine {
	// Set gin mode.
	gin.SetMode(Conf.Core.Mode)

	// Create the Gin engine.
	g := gin.New()

	// Routes
	router.Load(
		// Cores
		g,
		// Middlwares
		middleware.VersionMiddleware(),
		// middleware.Logging(),
		middleware.RequestId(),
	)
	return g
}

func autoTLSServer() *http.Server {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(Conf.Core.AutoTLS.Host),
		Cache:      autocert.DirCache(Conf.Core.AutoTLS.Folder),
	}
	return &http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   New(),
	}
}

func defaultTLSServer() *http.Server {
	return &http.Server{
		Addr:    "0.0.0.0:" + Conf.Core.TLS.Port,
		Handler: New(),
	}
}

func defaultServer() *http.Server {
	return &http.Server{
		Addr:    "0.0.0.0:" + Conf.Core.Port,
		Handler: New(),
	}
}

// RunHTTPServer provide run http or https protocol.
func RunHTTPServer() (err error) {
	if !Conf.Core.Enabled {
		log.Debug("httpd server is disabled.")
		return nil
	}

	if Conf.Core.AutoTLS.Enabled {
		s := autoTLSServer()
		handleSignal(s)
		log.Infof("1. Start to listening the incoming requests on https address")
		err = s.ListenAndServeTLS("", "")
	} else if Conf.Core.TLS.CertPath != "" && Conf.Core.TLS.KeyPath != "" {
		s := defaultTLSServer()
		handleSignal(s)
		log.Infof("2. Start to listening the incoming requests on https address: %s", Conf.Core.TLS.Port)
		err = s.ListenAndServeTLS(Conf.Core.TLS.CertPath, Conf.Core.TLS.KeyPath)
	} else {
		s := defaultServer()
		handleSignal(s)
		log.Infof("3. Start to listening the incoming requests on http address: %s", Conf.Core.Port)
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
			log.Errorf(err, "server close failed ")
		}

		service.DB.Close()

		log.Infof("apiserver exited")
		os.Exit(0)
	}()
}

// PingServer
func PingServer() (err error) {
	maxPingConf := Conf.Core.MaxPingCount
	for i := 0; i < maxPingConf; i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get("http://localhost:" + Conf.Core.Port + "/sd/health")
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
