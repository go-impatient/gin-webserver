package server

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

// WebServer 项目
type WebServer struct {
	config *config.Config
	// cache *storer.CacheStore
	// ...
}

var (
	Gorm           	*gorm.DB
	db   						*storer.Database
	Mail            *util.SendMail
)

func NewWebServer(cfg *config.Config) *WebServer {
	return &WebServer{
		config: cfg,
	}
}

// Init returns a app instance
func Init(w *WebServer) *gin.Engine {
	// Set gin mode.
	gin.SetMode(w.config.Core.Mode)

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

// Init initializes mail pkg.
func InitMail(w *WebServer) {
	Mail = util.SendMailNew(&util.SendMail{
		Enable: w.config.Mail.Enable,
		Smtp: w.config.Mail.Smtp,
		Port: w.config.Mail.Port,
		User: w.config.Mail.Username,
		Pass: w.config.Mail.Password,
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


func autoTLSServer(w *WebServer) *http.Server {
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(w.config.Core.AutoTLS.Host),
		Cache:      autocert.DirCache(w.config.Core.AutoTLS.Folder),
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
