package src

import (
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/moocss/go-webserver/src/config"
	"github.com/moocss/go-webserver/src/log"
	"github.com/moocss/go-webserver/src/router"
	"github.com/moocss/go-webserver/src/router/middleware"
	"github.com/moocss/go-webserver/src/storer"
	"github.com/moocss/go-webserver/src/util"
	"github.com/sevennt/wzap"
	"golang.org/x/crypto/acme/autocert"
)

var (
	A								*App
	Gorm           	*gorm.DB
	DB   						*storer.Database
	Mail            *util.SendMail
)

// App 项目
type App struct {
	config *config.Config
	serve  *gin.Engine
	// cache *storer.CacheStore
	// ...
}

func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
		serve: gin.New(),
	}
}

// Init initializes mail pkg.
func (app *App)InitMail() {
	Mail = util.SendMailNew(&util.SendMail{
		Enabled: app.config.Mail.Enable,
		Smtp: app.config.Mail.Smtp,
		Port: app.config.Mail.Port,
		Username: app.config.Mail.Username,
		Password: app.config.Mail.Password,
	})
}

// Init initializes log pkg.
func (app *App)InitLog() {
	wzap.SetDefaultDir("./log/")
	logger := wzap.New(
		wzap.WithOutput(
			wzap.WithLevelCombo(app.config.Log.Zap.Level),
			wzap.WithPath(app.config.Log.Zap.Path),
		),
		wzap.WithOutput(
			wzap.WithLevelCombo(app.config.Log.Console.Level),
			wzap.WithColorful(app.config.Log.Console.Color),
			wzap.WithPrefix(app.config.Log.Console.Prefix),
		),
	)
	wzap.SetDefaultLogger(logger)
}

// RunHTTPServer provide run http or https protocol.
func (app *App)RunHTTPServer() (err error) {
	if !app.config.Core.Enabled {
		log.Debug("httpd server is disabled.")
		return nil
	}

	if app.config.Core.AutoTLS.Enabled {
		s := autoTLSServer(app)
		handleSignal(s)
		log.Infof("1. Start to listening the incoming requests on https address")
		err = s.ListenAndServeTLS("", "")
	} else if app.config.Core.TLS.CertPath != "" && app.config.Core.TLS.KeyPath != "" {
		s := defaultTLSServer(app)
		handleSignal(s)
		log.Infof("2. Start to listening the incoming requests on https address: %s", app.config.Core.TLS.Port)
		err = s.ListenAndServeTLS(app.config.Core.TLS.CertPath, app.config.Core.TLS.KeyPath)
	} else {
		s := defaultServer(app)
		handleSignal(s)
		log.Infof("3. Start to listening the incoming requests on http address: %s", app.config.Core.Port)
		err = s.ListenAndServe()
	}

	return
}

func autoTLSServer(app *App) *http.Server {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(app.config.Core.AutoTLS.Host),
		Cache:      autocert.DirCache(app.config.Core.AutoTLS.Folder),
	}
	return &http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   serve(app),
	}
}

func defaultTLSServer(app *App) *http.Server {
	return &http.Server{
		Addr:    "0.0.0.0:" + app.config.Core.TLS.Port,
		Handler: serve(app),
	}
}

func defaultServer(app *App) *http.Server {
	return &http.Server{
		Addr:    "0.0.0.0:" + app.config.Core.Port,
		Handler: serve(app),
	}
}

// serve returns a app instance
func serve(app *App) *gin.Engine {
	// Set gin mode.
	gin.SetMode(app.config.Core.Mode)

	// Routes
	router.Load(
		// Cores
		app.serve,
		// Middlwares
		middleware.RequestId(),
	)
	return app.serve
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
func (app *App)PingServer() (err error) {
	maxPingConf := app.config.Core.MaxPingCount
	for i := 0; i < maxPingConf; i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get("http://localhost:" + app.config.Core.Port + "/sd/health")
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
