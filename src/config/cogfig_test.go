package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Test file is missing
func TestMissingFile(t *testing.T) {
	filePath := "test"
	_, err := Init(filePath)

	assert.NotNil(t, err)
}

type ConfigTestSuite struct {
	suite.Suite
	ConfDefault Config
	Conf        Config
}

func (suite *ConfigTestSuite) SetupTest() {
	var err error
	suite.ConfDefault, err = LoadConfig("")
	if err != nil {
		panic("failed to load default config.yml")
	}
	suite.Conf, err = LoadConfig("src/config.yml")
	if err != nil {
		panic("failed to load config.yml from file")
	}
}

func (suite *ConfigTestSuite) TestValidateConfDefault() {
	// Core
	assert.Equal(suite.T(), "apiserver", suite.ConfDefault.Core.Name)
	assert.Equal(suite.T(), true, suite.ConfDefault.Core.Enabled)
	assert.Equal(suite.T(), "", suite.ConfDefault.Core.Host)
	assert.Equal(suite.T(), "9090", suite.ConfDefault.Core.Port)
	assert.Equal(suite.T(), "debug", suite.ConfDefault.Core.Mode)
	assert.Equal(suite.T(), 2, suite.ConfDefault.Core.MaxPingCount)
	assert.Equal(suite.T(), "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", suite.ConfDefault.Core.JwtSecret)
	assert.Equal(suite.T(), "9098", suite.ConfDefault.Core.TLS.Port)
	assert.Equal(suite.T(), "src/server.crt", suite.ConfDefault.Core.TLS.CertPath)
	assert.Equal(suite.T(), "src/server.key", suite.ConfDefault.Core.TLS.KeyPath)

	// Log
	assert.Equal(suite.T(), true , suite.ConfDefault.Log.Console.Color)
	assert.Equal(suite.T(), "[webserver]", suite.ConfDefault.Log.Console.Prefix)
	assert.Equal(suite.T(), "debug", suite.ConfDefault.Log.Console.Level)
	assert.Equal(suite.T(), "webserver-api.log", suite.ConfDefault.Log.Zap.Path)
	assert.Equal(suite.T(), "debug", suite.ConfDefault.Log.Zap.Level)

	// Db
	assert.Equal(suite.T(), "db_apiserver", suite.ConfDefault.Db.DbName)
	assert.Equal(suite.T(), "127.0.0.1", suite.ConfDefault.Db.Host)
	assert.Equal(suite.T(), "3306", suite.ConfDefault.Db.Port)
	assert.Equal(suite.T(), "root", suite.ConfDefault.Db.Username)
	assert.Equal(suite.T(), "root", suite.ConfDefault.Db.Password)

	// Mail
	assert.Equal(suite.T(), true, suite.ConfDefault.Mail.Enable)
	assert.Equal(suite.T(), "smtp.exmail.qq.com", suite.ConfDefault.Mail.Smtp)
	assert.Equal(suite.T(), 465, suite.ConfDefault.Mail.Port)
	assert.Equal(suite.T(), "moocss@163.com", suite.ConfDefault.Mail.Username)
	assert.Equal(suite.T(), "", suite.ConfDefault.Mail.Password)
}

func (suite *ConfigTestSuite) TestValidateConf() {
	// Core
	assert.Equal(suite.T(), "", suite.Conf.Core.Host)
	assert.Equal(suite.T(), "9090", suite.Conf.Core.Port)
	assert.Equal(suite.T(), "debug", suite.Conf.Core.Mode)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
