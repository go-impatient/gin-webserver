package src

import (
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moocss/go-webserver/src/config"
	"github.com/moocss/go-webserver/src/log"
	"github.com/moocss/go-webserver/src/router"
	"github.com/moocss/go-webserver/src/router/middleware"
	"github.com/moocss/go-webserver/src/storer"
	"golang.org/x/crypto/acme/autocert"
	"github.com/jinzhu/gorm"
	"github.com/moocss/go-webserver/src/util"
	"github.com/sevennt/wzap"
	"github.com/spf13/viper"
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

// Init returns a app instance
func InitRouter(app *App) *gin.Engine {
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

// Init initializes mail pkg.
func InitMail(app *App) {
	Mail = util.SendMailNew(&util.SendMail{
		Enable: app.config.Mail.Enable,
		Smtp: app.config.Mail.Smtp,
		Port: app.config.Mail.Port,
		User: app.config.Mail.Username,
		Pass: app.config.Mail.Password,
	})
}

// Init initializes log pkg.
func InitLog() {
	logger := wzap.New(
		wzap.WithOutputKV(viper.GetStringMap("logger.console")),
		wzap.WithOutputKV(viper.GetStringMap("logger.zap")),
	)
	wzap.SetDefaultLogger(logger)
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
		Handler:   Init(w),
	}
}

func defaultTLSServer(w *WebServer) *http.Server {
	return &http.Server{
		Addr:    "0.0.0.0:" + w.config.Core.TLS.Port,
		Handler: Init(w),
	}
}

func defaultServer(w *WebServer) *http.Server {
	return &http.Server{
		Addr:    "0.0.0.0:" + w.config.Core.Port,
		Handler: Init(w),
	}
}

// RunHTTPServer provide run http or https protocol.
func RunHTTPServer(w *WebServer) (err error) {
	if !w.config.Core.Enabled {
		log.Debug("httpd server is disabled.")
		return nil
	}

	if w.config.Core.AutoTLS.Enabled {
		s := autoTLSServer(w)
		handleSignal(s)
		log.Infof("1. Start to listening the incoming requests on https address")
		err = s.ListenAndServeTLS("", "")
	} else if w.config.Core.TLS.CertPath != "" && w.config.Core.TLS.KeyPath != "" {
		s := defaultTLSServer(w)
		handleSignal(s)
		log.Infof("2. Start to listening the incoming requests on https address: %s", w.config.Core.TLS.Port)
		err = s.ListenAndServeTLS(w.config.Core.TLS.CertPath, w.config.Core.TLS.KeyPath)
	} else {
		s := defaultServer(w)
		handleSignal(s)
		log.Infof("3. Start to listening the incoming requests on http address: %s", w.config.Core.Port)
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
func PingServer(w *WebServer) (err error) {
	maxPingConf := w.config.Core.MaxPingCount
	for i := 0; i < maxPingConf; i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get("http://localhost:" + w.config.Core.Port + "/sd/health")
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
