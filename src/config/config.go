package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/moocss/go-webserver/src/log"
	"github.com/spf13/viper"
)

var defaultConfig = []byte(`
core:
  enabled: true                   # enabale httpd server
  mode: "debug"                   # 开发模式, debug, release, test
  name: "apiserver"               # API Server的名字
  host: ""                        # ip address to bind (default: any)
  port: "9090"                    # HTTP 绑定端口.
  max_ping_count: 2               # pingServer函数try的次数
  jwt_secret: "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5"
  tls:
    port: "9098"
    cert_path: ""                 # src/config/server.crt
    key_path: ""                  # src/config/server.key
  auto_tls:
    enabled: false                # Automatically install TLS certificates from Let's Encrypt.
    folder: ".cache"              # folder for storing TLS certificates
    host: ""                      # which domains the Let's Encrypt will attempt

log:
    console:
      color: true
      prefix: "[webserver]"
      level: "debug"
    zap:
      path: "webserver-api.log"
      level: "debug"

db:
  db_name: "db_apiserver"
  host: "127.0.0.1"
  port: "3306"
  username: "root"
  password: "123456"
  charset: "utf8mb4"
  unix: ""
  table_prefix: ""
  max_idle_conns: ""
  max_open_conns: ""
  conn_max_lift_time: ""

mail:
  enabled: true                    # 是否开启邮箱发送功能
  smtp_host: "smtp.exmail.qq.com"  # 邮件smtp地址
  smtp_port: 465
  smtp_username: "moocss@163.com"
  smtp_password: ""
`)

type (
	Config struct {
		Core 	*ConfigCore `yaml:"core"`
		Log  	*ConfigLog  `yaml:"log"`
		Db   	*ConfigDb   `yaml:"db"`
		Mail 	*ConfigMail `yaml:"mail"`
	}
	// ConfigCore is sub section of config.
	ConfigCore struct {
		Enabled      bool           `yaml:"enabled"`
		Mode         string         `yaml:"mode"`
		Name         string         `yaml:"name"`
		Host         string         `yaml:"host"`
		Port         string         `yaml:"port"`
		MaxPingCount int            `yaml:"max_ping_count"`
		JwtSecret    string         `yaml:"jwt_secret"`
		TLS          ConfigTLS     	`yaml:"tls"`
		AutoTLS      ConfigAutoTLS 	`yaml:"auto_tls"`
	}

	// ConfigTLS support tls
	ConfigTLS struct {
		Port     string `yaml:"port"`
		CertPath string `yaml:"cert_path"`
		KeyPath  string `yaml:"key_path"`
	}

	// ConfigAutoTLS support Let's Encrypt setting.
	ConfigAutoTLS struct {
		Enabled bool   `yaml:"enabled"`
		Folder  string `yaml:"folder"`
		Host    string `yaml:"host"`
	}

	// ConfigLog is sub section of config.
	ConfigLog struct {
		Console  ConfigLogConsole 	`yaml:"console"`
		Zap			 ConfigLogZap 			`yaml:"zap"`
	}
	ConfigLogConsole struct {
		Color 	bool   `yaml:"color"`
		Prefix 	string `yaml:"prefix"`
		Level 	string `yaml:"level"`
	}
	ConfigLogZap struct {
		Path  	string `yaml:"path"`
		Level 	string `yaml:"level"`
	}

	// ConfigDb is sub section of config.
	ConfigDb struct {
		Unix            string `yaml:"unix"`
		Host            string `yaml:"host"`
		Port            string `yaml:"port"`
		Charset         string `yaml:"charset"`
		DbName          string `yaml:"db_name"`
		Username        string `yaml:"username"`
		Password        string `yaml:"password"`
		TablePrefix     string `yaml:"table_prefix"`
		MaxIdleConns    int    `yaml:"max_idle_conns"`
		MaxOpenConns    int    `yaml:"max_open_conns"`
		ConnMaxLifeTime int    `yaml:"conn_max_lift_time"`
	}

	// ConfigMail is sub section of config
	ConfigMail struct {
		Enable   bool     `yaml:"enable"`
		Smtp     string   `yaml:"smtp_host"`
		Port     int  		`yaml:"smtp_port"`
		Username string   `yaml:"smtp_username"`
		Password string   `yaml:"smtp_password"`
	}
)

// 加载配置文件
func LoadConfig(confPath string) (Config, error) {
	var cfg Config

	// 设置配置文件格式为YAML
	viper.SetConfigType("yaml")

	// 读取匹配的环境变量
	viper.AutomaticEnv()

	// 读取环境变量的前缀为bear
	viper.SetEnvPrefix("bear")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if confPath != "" {
		// 如果指定了配置文件路径，则解析指定的配置文件路径
		viper.SetConfigFile(confPath)
	} else {
		// 如果没有指定配置文件，则解析默认的配置文件
		// Search config in home directory with name ".webserver" (without extension).
		viper.AddConfigPath("/etc/webserver/")
		viper.AddConfigPath("$HOME/.webserver")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("解析配置文件失败:", viper.ConfigFileUsed())
		return cfg, err
	} else {
		// load default config
		err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
		if err != nil {
			log.Fatal("读取默认配置失败: " + err.Error())
			return cfg, err
		}
	}

	// 将新配置解组到我们的运行时配置结构中。
	//if err := viper.Unmarshal(Bear.C); err != nil {
	//	log.Fatal("解密配置失败: " + err.Error())
	//	return err
	//}

	// Core
	cfg.Core.Enabled = viper.GetBool("core.enabled")
	cfg.Core.Mode = viper.GetString("core.mode")
	cfg.Core.Name = viper.GetString("core.name")
	cfg.Core.Host = viper.GetString("core.host")
	cfg.Core.Port = viper.GetString("core.port")
	cfg.Core.MaxPingCount = viper.GetInt("core.max_ping_count")
	cfg.Core.JwtSecret = viper.GetString("core.jwt_secret")
	cfg.Core.TLS.Port = viper.GetString("core.tls.port")
	cfg.Core.TLS.CertPath = viper.GetString("core.tls.cert_path")
	cfg.Core.TLS.KeyPath = viper.GetString("core.tls.key_path")
	cfg.Core.AutoTLS.Enabled = viper.GetBool("core.auto_tls.enabled")
	cfg.Core.AutoTLS.Folder = viper.GetString("core.auto_tls.folder")
	cfg.Core.AutoTLS.Host = viper.GetString("core.auto_tls.host")

	// Log
	cfg.Log.Console.Color = viper.GetBool("log.console.color")
	cfg.Log.Console.Prefix = viper.GetString("log.console.prefix")
	cfg.Log.Console.Level = viper.GetString("log.console.level")
	cfg.Log.Zap.Path = viper.GetString("log.zap.path")
	cfg.Log.Zap.Level = viper.GetString("log.zap.level")

	// Db
	cfg.Db.Unix = viper.GetString("db.unix")
	cfg.Db.Host = viper.GetString("db.host")
	cfg.Db.Port = viper.GetString("db.port")
	cfg.Db.Charset = viper.GetString("db.charset")
	cfg.Db.DbName = viper.GetString("db.db_name")
	cfg.Db.Username = viper.GetString("db.username")
	cfg.Db.Password = viper.GetString("db.password")
	cfg.Db.TablePrefix = viper.GetString("db.table_prefix")
	cfg.Db.MaxIdleConns = viper.GetInt("max_idle_conns")
	cfg.Db.MaxOpenConns = viper.GetInt("max_open_conns")
	cfg.Db.ConnMaxLifeTime = viper.GetInt("conn_max_lift_time")

	// Mail
	cfg.Mail.Enable = viper.GetBool("mail.enable")
	cfg.Mail.Smtp = viper.GetString("mail.smtp_host")
	cfg.Mail.Port = viper.GetInt("mail.smtp_port")
	cfg.Mail.Username = viper.GetString("mail.smtp_username")
	cfg.Mail.Password = viper.GetString("mail.smtp_password")

	return cfg, nil
}

// 监控配置文件变化并热加载程序
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}
