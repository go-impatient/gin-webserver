package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/moocss/go-webserver/src/log"
	"github.com/moocss/go-webserver/src/storer"
)

var defaultConf = []byte(`
core:
  enabled: true                   # enabale httpd server
  mode: "debug"                   # 开发模式, debug, release, test
  name: "apiserver"               # API Server的名字
  address: ""                     # ip address to bind (default: any)
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
db:
  name: "db_apiserver"
  addr: "127.0.0.1:3306"
  username: "root"
  password: "123456"

docker_db:
  name: "db_apiserver"
  addr: "127.0.0.1:3306"
  username: "root"
  password: "123456"
`)

type Config struct {
	Core     SectionCore     `yaml:"core"`
	Log      SectionLog      `yaml:"log"`
	Db       SectionDb       `yaml:"db"`
	DockerDb SectionDockerDb `yaml:"db"`
}

// SectionCore is sub section of config.
type SectionCore struct {
	Enabled      bool           `yaml:"enabled"`
	Mode         string         `yaml:"mode"`
	Name         string         `yaml:"name"`
	Host         string         `yaml:"host"`
	Port         string         `yaml:"port"`
	MaxPingCount int            `yaml:"max_ping_count"`
	JwtSecret    string         `yaml:"jwt_secret"`
	TLS          SectionTLS     `yaml:"tls"`
	AutoTLS      SectionAutoTLS `yaml:"auto_tls"`
}

// SectionTLS support tls
type SectionTLS struct {
	Port     string `yaml:"port"`
	CertPath string `yaml:"cert_path"`
	KeyPath  string `yaml:"key_path"`
}

// SectionAutoTLS support Let's Encrypt setting.
type SectionAutoTLS struct {
	Enabled bool   `yaml:"enabled"`
	Folder  string `yaml:"folder"`
	Host    string `yaml:"host"`
}

// SectionLog is sub section of config.
type SectionLog struct {
	Writers        string `yaml:"writers"`
	LoggerLevel    string `yaml:"logger_level"`
	LoggerFile     string `yaml:"logger_file"`
	LogFormatText  bool   `yaml:"log_format_text"`
	RollingPolicy  string `yaml:"rollingPolicy"`
	LogRotateDate  int    `yaml:"log_rotate_date"`
	LogRotateSize  int    `yaml:"log_rotate_size"`
	LogBackupCount int    `yaml:"log_backup_count"`
}

// SectionDb is sub section of config.
type SectionDb struct {
	Name     string `yaml:"name"`
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// SectionDb is sub section of config.
type SectionDockerDb struct {
	Name     string `yaml:"name"`
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Init(confPath string) error {

	// 初始化配置文件
	if err := initConfig(confPath); err != nil {
		return err
	}

	// 监控配置文件变化并热加载程序
	watchConfig()

	return nil
}

// 加载配置文件
func initConfig(confPath string) error {
	var conf Config

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
		// Search config in home directory with name ".bear" (without extension).
		viper.AddConfigPath("/etc/bear/")
		viper.AddConfigPath("$HOME/.bear")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		return err
	} else {
		// load default config
		err := viper.ReadConfig(bytes.NewBuffer(defaultConf))
		if err != nil {
			return err
			log.Fatal("读取默认失败: " + err.Error())
		}
	}

	// 将新配置解组到我们的运行时配置结构中。
	if err := viper.Unmarshal(Bear.C); err != nil {
		log.Fatal("解密配置失败: " + err.Error())
		return err
	}

	// 监控配置文件变化并热加载程序
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})

	// Core
	conf.Core.Enabled = viper.GetBool("core.enabled")
	conf.Core.Mode = viper.GetString("core.mode")
	conf.Core.Name = viper.GetString("core.name")
	conf.Core.Host = viper.GetString("core.host")
	conf.Core.Port = viper.GetString("core.port")
	conf.Core.MaxPingCount = viper.GetInt("core.max_ping_count")
	conf.Core.JwtSecret = viper.GetString("core.jwt_secret")
	conf.Core.TLS.Port = viper.GetString("core.tls.port")
	conf.Core.TLS.CertPath = viper.GetString("core.tls.cert_path")
	conf.Core.TLS.KeyPath = viper.GetString("core.tls.key_path")
	conf.Core.AutoTLS.Enabled = viper.GetBool("core.auto_tls.enabled")
	conf.Core.AutoTLS.Folder = viper.GetString("core.auto_tls.folder")
	conf.Core.AutoTLS.Host = viper.GetString("core.auto_tls.host")

	// Db
	conf.Db.Name = viper.GetString("db.name")
	conf.Db.Addr = viper.GetString("db.addr")
	conf.Db.Username = viper.GetString("db.username")
	conf.Db.Password = viper.GetString("db.password")

	// DockerDb
	conf.DockerDb.Name = viper.GetString("docker_db.name")
	conf.DockerDb.Addr = viper.GetString("docker_db.addr")
	conf.DockerDb.Username = viper.GetString("docker_db.username")
	conf.DockerDb.Password = viper.GetString("docker_db.password")

	return nil
}

// 监控配置文件变化并热加载程序
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}

// 项目
type bear struct {
	C     *Config
	Cache *storer.CacheStore
	// ...
}

// Bear 包含全局信息，更重要是配置信息
var Bear = bear{
	C: &Config{},
	// ...
}
