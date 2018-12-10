package log

import (
	"github.com/sevennt/wzap"
	"github.com/spf13/viper"
)

// Init initializes log pkg.
func Init() {
	logger := wzap.New(
		wzap.WithOutputKV(viper.GetStringMap("logger.console")),
		wzap.WithOutputKV(viper.GetStringMap("logger.zap")),
	)
	wzap.SetDefaultLogger(logger)
}
