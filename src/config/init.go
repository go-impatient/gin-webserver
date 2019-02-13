package config

// Init initializes config pkg.
func Init(confPath string) (Config, error) {
	var cfg Config
	// 初始化配置文件
	if cfg, err := LoadConfig(confPath); err != nil {
		return cfg, err
	}

	// 监控配置文件变化并热加载程序
	watchConfig()

	return cfg, nil
}
