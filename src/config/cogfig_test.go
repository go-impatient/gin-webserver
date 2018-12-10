package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Test file is missing
func TestMissingFile(t *testing.T) {
	filename := "test"
	_, err := Init(filename)

	assert.NotNil(t, err)
}

type ConfigTestSuite struct {
	suite.Suite
	ConfDefault 	ConfYaml
	Conf        	ConfYaml
}

func (suite *ConfigTestSuite) SetupTest() {
	var err error
	suite.ConfDefault, err = Init("")
	if err != nil {
		panic("failed to load default config.yml")
	}
	suite.Conf, err = Init("src/config/config")
	if err != nil {
		panic("failed to load config.yml from file")
	}
}

func (suite *ConfigTestSuite) TestValidateConfDefault() {
	// Core
	assert.Equal(suite.T(), "apiserver", suite.ConfDefault.Core.Name)
	assert.Equal(suite.T(), true, suite.ConfDefault.Core.Enabled)
	assert.Equal(suite.T(), "", suite.ConfDefault.Core.Address)
	assert.Equal(suite.T(), "9090", suite.ConfDefault.Core.Port)
	assert.Equal(suite.T(), "debug", suite.ConfDefault.Core.Mode)
	assert.Equal(suite.T(), 2, suite.ConfDefault.Core.MaxPingCount)
	assert.Equal(suite.T(), "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", suite.ConfDefault.Core.JwtSecret)
	assert.Equal(suite.T(), "9098", suite.ConfDefault.Core.TLS.Port)
	assert.Equal(suite.T(), "src/config/server.crt", suite.ConfDefault.Core.TLS.CertPath)
	assert.Equal(suite.T(), "src/config/server.key", suite.ConfDefault.Core.TLS.KeyPath)

	// Log
	assert.Equal(suite.T(), "file,stdout", suite.ConfDefault.Log.Writers)
	assert.Equal(suite.T(), "DEBUG", suite.ConfDefault.Log.LoggerLevel)
	assert.Equal(suite.T(), "log/apiserver.log", suite.ConfDefault.Log.LoggerFile)
	assert.Equal(suite.T(), false, suite.ConfDefault.Log.LogFormatText)
	assert.Equal(suite.T(), "size", suite.ConfDefault.Log.RollingPolicy)
	assert.Equal(suite.T(), 1, suite.ConfDefault.Log.LogRotateDate)
	assert.Equal(suite.T(), 1, suite.ConfDefault.Log.LogRotateSize)
	assert.Equal(suite.T(), 7, suite.ConfDefault.Log.LogBackupCount)

	// Db
	assert.Equal(suite.T(), "db_apiserver", suite.ConfDefault.Db.Name)
	assert.Equal(suite.T(), "127.0.0.1:3306", suite.ConfDefault.Db.Addr)
	assert.Equal(suite.T(), "root", suite.ConfDefault.Db.Username)
	assert.Equal(suite.T(), "root", suite.ConfDefault.Db.Password)

}

func (suite *ConfigTestSuite) TestValidateConf() {
	// Core
	assert.Equal(suite.T(), "", suite.Conf.Core.Address)
	assert.Equal(suite.T(), "9090", suite.Conf.Core.Port)
	assert.Equal(suite.T(), "debug", suite.Conf.Core.Mode)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
