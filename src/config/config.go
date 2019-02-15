package config

import (
	"bytes"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/moocss/go-webserver/src/log"
	"github.com/spf13/viper"
)

var defaultConfig = []byte(`
core:
  enabled: true                   # enabale httpd server
  mode: "dev"             				# dev(debug), prod(release), test
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
  dialect: "mysql"
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

cache:
  type: "none"
  timeout: 60
  redis:
    host: "localhost"
    port: 6379
    password: ""
    db: 0
    keyprefix: "__:::webserver:"
`)

type (
	Config struct {
		Core  *ConfigCore  `yaml:"core"`
		Log   *ConfigLog   `yaml:"log"`
		Db    *ConfigDb    `yaml:"db"`
		Mail  *ConfigMail  `yaml:"mail"`
		Cache *ConfigCache `yaml:"cache"`
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
		TLS          *ConfigTLS     `yaml:"tls"`
		AutoTLS      *ConfigAutoTLS `yaml:"auto_tls"`
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
		Console *ConfigLogConsole `yaml:"console"`
		Zap     *ConfigLogZap     `yaml:"zap"`
	}
	ConfigLogConsole struct {
		Color  bool   `yaml:"color"`
		Prefix string `yaml:"prefix"`
		Level  string `yaml:"level"`
	}
	ConfigLogZap struct {
		Path  string `yaml:"path"`
		Level string `yaml:"level"`
	}

	// ConfigDb is sub section of config.
	ConfigDb struct {
		Unix            string `yaml:"unix"`
		Host            string `yaml:"host"`
		Port            string `yaml:"port"`
		Charset         string `yaml:"charset"`
		Dialect         string `yaml:"dialect"`
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
		Enable   bool   `yaml:"enable"`
		Smtp     string `yaml:"smtp_host"`
		Port     int    `yaml:"smtp_port"`
		Username string `yaml:"smtp_username"`
		Password string `yaml:"smtp_password"`
	}

	ConfigCache struct {
		Type    string            `yaml:"type"`
		Timeout int32             `yaml:"timeout"`
		Redis   *ConfigCacheRedis `yaml:"redis"`
	}

	ConfigCacheRedis struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		Password  string `yaml:"password"`
		DB        int    `yaml:"db"`
		KeyPrefix string `yaml:"keyprefix"`
	}
)

// Init initializes config pkg.
func Init(confPath string) (*Config, error) {
	// 初始化配置文件
	cfg, err := LoadConfig(confPath)
	if err != nil {
		return nil, err
	}

	// 监控配置文件变化并热加载程序
	watchConfig()

	return cfg, nil
}

// 加载配置文件
func LoadConfig(confPath string) (*Config, error) {
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
		log.Errorf("解析配置文件失败: %s", viper.ConfigFileUsed())
		return nil, err
	} else {
		// load default config
		err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
		if err != nil {
			log.Errorf("读取默认配置失败: %s", err.Error())
			return nil, err
		}
	}

	cfg := &Config{
		Core: &ConfigCore{
			Enabled:      viper.GetBool("core.enabled"),
			Mode:         viper.GetString("core.mode"),
			Name:         viper.GetString("core.name"),
			Host:         viper.GetString("core.host"),
			Port:         viper.GetString("core.port"),
			MaxPingCount: viper.GetInt("core.max_ping_count"),
			JwtSecret:    viper.GetString("core.jwt_secret"),
			TLS: &ConfigTLS{
				Port:     viper.GetString("core.tls.port"),
				CertPath: viper.GetString("core.tls.cert_path"),
				KeyPath:  viper.GetString("core.tls.key_path"),
			},
			AutoTLS: &ConfigAutoTLS{
				Enabled: viper.GetBool("core.auto_tls.enabled"),
				Folder:  viper.GetString("core.auto_tls.folder"),
				Host:    viper.GetString("core.auto_tls.host"),
			},
		},
		Log: &ConfigLog{
			Console: &ConfigLogConsole{
				Color:  viper.GetBool("log.console.color"),
				Prefix: viper.GetString("log.console.prefix"),
				Level:  viper.GetString("log.console.level"),
			},
			Zap: &ConfigLogZap{
				Path:  viper.GetString("log.zap.path"),
				Level: viper.GetString("log.zap.level"),
			},
		},
		Db: &ConfigDb{
			Unix:            viper.GetString("db.unix"),
			Host:            viper.GetString("db.host"),
			Port:            viper.GetString("db.port"),
			Charset:         viper.GetString("db.charset"),
			DbName:          viper.GetString("db.db_name"),
			Username:        viper.GetString("db.username"),
			Password:        viper.GetString("db.password"),
			TablePrefix:     viper.GetString("db.table_prefix"),
			MaxIdleConns:    viper.GetInt("db.max_idle_conns"),
			MaxOpenConns:    viper.GetInt("db.max_open_conns"),
			ConnMaxLifeTime: viper.GetInt("db.conn_max_lift_time"),
		},
		Mail: &ConfigMail{
			Enable:   viper.GetBool("mail.enable"),
			Smtp:     viper.GetString("mail.smtp_host"),
			Port:     viper.GetInt("mail.smtp_port"),
			Username: viper.GetString("mail.smtp_username"),
			Password: viper.GetString("mail.smtp_password"),
		},
		Cache: &ConfigCache{
			Type:    viper.GetString("cache.type"),
			Timeout: viper.GetInt32("cache.timeout"),
			Redis: &ConfigCacheRedis{
				Host:      viper.GetString("cache.redis.host"),
				Port:      viper.GetInt("cache.redis.host"),
				Password:  viper.GetString("cache.redis.password"),
				DB:        viper.GetInt("cache.redis.db"),
				KeyPrefix: viper.GetString("cache.redis.keyprefix"),
			},
		},
	}

	return cfg, nil
}

// 监控配置文件变化并热加载程序
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}
