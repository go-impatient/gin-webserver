package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/moocss/go-webserver/src/log"
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
	Name string
}

type ConfYaml struct {
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
	Host      	 string         `yaml:"host"`
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

func Init(configName string) (ConfYaml, error) {
	c := Config{
		Name: configName,
	}

	// 初始化配置文件
	conf, err := c.initConfig()
	if err != nil {
		return conf, nil
	}

	// 监控配置文件变化并热加载程序
	c.watchConfig()

	return conf, nil
}

// 初始化配置文件
func (c *Config) initConfig() (ConfYaml, error) {
	var conf ConfYaml

	// 设置配置文件格式为YAML
	viper.SetConfigType("yaml")

	// 读取匹配的环境变量
	viper.AutomaticEnv()

	// 读取环境变量的前缀为BEARSERVER
	viper.SetEnvPrefix("BEARSERVER")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if c.Name != "" {
		// 如果指定了配置文件，则解析指定的配置文件
		viper.SetConfigFile(c.Name)
	} else {
		// 如果没有指定配置文件，则解析默认的配置文件
		viper.AddConfigPath("src/config")
		viper.SetConfigName("config")
	}

	// viper解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// load default config
		viper.ReadConfig(bytes.NewBuffer(defaultConf))
	}

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

	return conf, nil
}

// 监控配置文件变化并热加载程序
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}
