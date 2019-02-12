package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/moocss/go-webserver/src/log"
	"github.com/spf13/viper"
	"github.com/moocss/go-webserver/src/server"
	"github.com/moocss/go-webserver/src/util"
)

var defaultConfig = []byte(``)

type (
	Config struct {
		Core *ConfigCore `yaml:"core"`
		Log  *ConfigLog  `yaml:"log"`
		Db   *ConfigDb   `yaml:"db"`
		Mail *ConfigMail `yaml:"mail"`
	}
	// ConfigCore is sub section of config.
	ConfigCore struct {
		Enabled      bool           `yaml:"enabled"`
		Mode         string         `yaml:"mode"`
		Name         string         `yaml:"name"`
		Host         string         `yaml:"host"`
		Port         string         `yaml:"port"`
		ReadTimeout  int            `yaml:"read_timeout"`
		WriteTimeout int            `yaml:"write_timeout"`
		IdleTimeout  int            `yaml:"idle_timeout"`
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

	// SectionAutoTLS support Let's Encrypt setting.
	ConfigAutoTLS struct {
		Enabled bool   `yaml:"enabled"`
		Folder  string `yaml:"folder"`
		Host    string `yaml:"host"`
	}

	// SectionLog is sub section of config.
	ConfigLog struct {
		Writers        string `yaml:"writers"`
		LoggerLevel    string `yaml:"logger_level"`
		LoggerFile     string `yaml:"logger_file"`
		LogFormatText  bool   `yaml:"log_format_text"`
		RollingPolicy  string `yaml:"rolling_policy"`
		LogRotateDate  int    `yaml:"log_rotate_date"`
		LogRotateSize  int    `yaml:"log_rotate_size"`
		LogBackupCount int    `yaml:"log_backup_count"`
	}

	// SectionDb is sub section of config.
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

	// SectionMail is sub section of config
	ConfigMail struct {
		Enable   int    `yaml:"enable"`
		Smtp     string `yaml:"smtp"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
)

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
	var conf *Config

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
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		return err
	} else {
		// load default config
		err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
		if err != nil {
			log.Fatal("读取默认失败: " + err.Error())
			return err
		}
	}

	// 将新配置解组到我们的运行时配置结构中。
	//if err := viper.Unmarshal(Bear.C); err != nil {
	//	log.Fatal("解密配置失败: " + err.Error())
	//	return err
	//}

	// 监控配置文件变化并热加载程序
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})

	//cfg := &Config{
	//
	//}

	return nil
}

// 监控配置文件变化并热加载程序
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}
