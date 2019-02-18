package src

import (
	"crypto/tls"
	"errors"
	"github.com/moocss/go-webserver/src/service"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sevennt/wzap"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"

	"github.com/moocss/go-webserver/src/config"
	"github.com/moocss/go-webserver/src/pkg/log"
	"github.com/moocss/go-webserver/src/pkg/mail"
	"github.com/moocss/go-webserver/src/router"
	"github.com/moocss/go-webserver/src/router/middleware"
	"github.com/moocss/go-webserver/src/storer"
	"github.com/moocss/go-webserver/src/util"
)

var (
	A          *App
	Orm        *gorm.DB
	DbInstance *storer.Database
	Mail       *mail.SendMail
)

// App 项目
type App struct {
	config *config.Config
	// cache *storer.CacheStore
	// ...
}

func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

func (app *App) InitDB() {
	DbInstance = storer.NewDatabase(app.config.Db)
	db, err := DbInstance.Open()
	if err != nil {
		panic(err)
	}

	Orm = DbInstance.Self

	defer db.Close()
}

// Init initializes mail pkg.
func (app *App) InitMail() {
	Mail = mail.SendMailNew(&mail.SendMail{
		Enabled:  app.config.Mail.Enable,
		Smtp:     app.config.Mail.Smtp,
		Port:     app.config.Mail.Port,
		Username: app.config.Mail.Username,
		Password: app.config.Mail.Password,
	})
}

// Init initializes log pkg.
func (app *App) InitLog() {
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
func (app *App) RunHTTPServer() (err error) {
	if !app.config.Core.Enabled {
		log.Debug("httpd server is disabled.")
		return nil
	}

	if app.config.Core.AutoTLS.Enabled {
		return autoTLSServer(app)
	} else if app.config.Core.TLS.CertPath != "" && app.config.Core.TLS.KeyPath != "" {
		return defaultTLSServer(app)
	} else {
		return defaultServer(app)
	}
}

func autoTLSServer(app *App) error {
	var g errgroup.Group

	dir := util.CacheDir()
	_ = os.MkdirAll(dir, 0700)

	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(app.config.Core.AutoTLS.Host),
		// Cache:      autocert.DirCache(app.config.Core.AutoTLS.Folder),
		Cache: autocert.DirCache(dir),
	}

	g.Go(func() error {
		return http.ListenAndServe(":http", manager.HTTPHandler(http.HandlerFunc(redirect)))
	})

	g.Go(func() error {
		serve := &http.Server{
			Addr: ":https",
			TLSConfig: &tls.Config{
				GetCertificate: manager.GetCertificate,
				NextProtos:     []string{"http/1.1"}, // disable h2 because Safari :(
			},
			Handler: serve(app),
		}
		handleSignal(serve)
		log.Info("Start to listening the incoming requests on https address")
		return serve.ListenAndServeTLS("", "")
	})

	return g.Wait()
}

func defaultTLSServer(app *App) error {
	var g errgroup.Group
	g.Go(func() error {
		return http.ListenAndServe(":http", http.HandlerFunc(redirect))
	})
	g.Go(func() error {
		serve := &http.Server{
			Addr:    ":https",
			Handler: serve(app),
			TLSConfig: &tls.Config{
				NextProtos: []string{"http/1.1"}, // disable h2 because Safari :(
			},
		}
		handleSignal(serve)
		log.Infof("Start to listening the incoming requests on https address: %s", app.config.Core.TLS.Port)
		return serve.ListenAndServeTLS(
			app.config.Core.TLS.CertPath,
			app.config.Core.TLS.KeyPath,
		)
	})
	return g.Wait()
}

func defaultServer(app *App) error {
	serve := &http.Server{
		Addr:    "0.0.0.0:" + app.config.Core.Port,
		Handler: serve(app),
	}

	handleSignal(serve)
	log.Infof("Start to listening the incoming requests on http address: %s", app.config.Core.Port)
	return serve.ListenAndServe()
}

// redirect ...
func redirect(w http.ResponseWriter, req *http.Request) {
	var serverHost string = A.config.Core.Host
	serverHost = strings.TrimPrefix(serverHost, "http://")
	serverHost = strings.TrimPrefix(serverHost, "https://")
	req.URL.Scheme = "https"
	req.URL.Host = serverHost

	w.Header().Set("Strict-Transport-Security", "max-age=31536000")

	http.Redirect(w, req, req.URL.String(), http.StatusMovedPermanently)
}

// serve returns a app instance
func serve(app *App) *gin.Engine {
	// Set gin mode.
	setRuntimeMode(app.config.Core.Mode)

	// Setup Business Layer
	s := service.NewService()

	// Setup the server
	handler := router.Load(
		// Services
		s,

		// Middlwares
		middleware.RequestId(),
	)

	return handler
}

func setRuntimeMode(mode string) {
	switch mode {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		panic("unknown mode")
	}
}

// handleSignal handles system signal for graceful shutdown.
func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		log.Infof("got signal [%s], exiting apiserver now", s)
		if err := server.Close(); nil != err {
			log.Error("server close failed ", err)
		}

		log.Info("apiserver exited")
		os.Exit(0)
	}()
}

// PingServer
func (app *App) PingServer() (err error) {
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
