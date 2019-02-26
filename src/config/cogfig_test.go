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
	ConfDefault *Config
	Conf        *Config
}

func (suite *ConfigTestSuite) SetupTest() {
	var err error
	suite.ConfDefault, err = Init("")
	if err != nil {
		panic("failed to load default config.yaml")
	}
	suite.Conf, err = Init("testdata/config.yaml")
	if err != nil {
		panic("failed to load config.yaml from file")
	}
}

func (suite *ConfigTestSuite) TestValidateConfDefault() {
	// Core
	assert.Equal(suite.T(), "webserver", suite.ConfDefault.Core.Name)
	assert.Equal(suite.T(), true, suite.ConfDefault.Core.Enabled)
	assert.Equal(suite.T(), "", suite.ConfDefault.Core.Host)
	assert.Equal(suite.T(), "9090", suite.ConfDefault.Core.Port)
	assert.Equal(suite.T(), "dev", suite.ConfDefault.Core.Mode)
	assert.Equal(suite.T(), 2, suite.ConfDefault.Core.MaxPingCount)
	assert.Equal(suite.T(), "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", suite.ConfDefault.Core.JwtSecret)
	assert.Equal(suite.T(), "9098", suite.ConfDefault.Core.TLS.Port)
	assert.Equal(suite.T(), "", suite.ConfDefault.Core.TLS.CertPath)
	assert.Equal(suite.T(), "", suite.ConfDefault.Core.TLS.KeyPath)

	// Log
	assert.Equal(suite.T(), true, suite.ConfDefault.Log.Console.Color)
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
	assert.Equal(suite.T(), true, suite.ConfDefault.Db.LogMode)

	// Mail
	assert.Equal(suite.T(), false, suite.ConfDefault.Mail.Enable)
	assert.Equal(suite.T(), "smtp.exmail.qq.com", suite.ConfDefault.Mail.Smtp)
	assert.Equal(suite.T(), 465, suite.ConfDefault.Mail.Port)
	assert.Equal(suite.T(), "moocss@163.com", suite.ConfDefault.Mail.Username)
	assert.Equal(suite.T(), "", suite.ConfDefault.Mail.Password)
}

func (suite *ConfigTestSuite) TestValidateConf() {
	// Core
	assert.Equal(suite.T(), "webserver", suite.Conf.Core.Name)
	assert.Equal(suite.T(), true, suite.Conf.Core.Enabled)
	assert.Equal(suite.T(), "", suite.Conf.Core.Host)
	assert.Equal(suite.T(), "9090", suite.Conf.Core.Port)
	assert.Equal(suite.T(), "dev", suite.Conf.Core.Mode)
	assert.Equal(suite.T(), 2, suite.Conf.Core.MaxPingCount)
	assert.Equal(suite.T(), "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", suite.Conf.Core.JwtSecret)
	assert.Equal(suite.T(), "9098", suite.Conf.Core.TLS.Port)
	assert.Equal(suite.T(), "", suite.Conf.Core.TLS.CertPath)
	assert.Equal(suite.T(), "", suite.Conf.Core.TLS.KeyPath)

	// Log
	assert.Equal(suite.T(), "./log/", suite.Conf.Log.DefaultDir)
	assert.Equal(suite.T(), true, suite.Conf.Log.Console.Color)
	assert.Equal(suite.T(), "[webserver]", suite.Conf.Log.Console.Prefix)
	assert.Equal(suite.T(), "debug", suite.Conf.Log.Console.Level)
	assert.Equal(suite.T(), "webserver-api.log", suite.Conf.Log.Zap.Path)
	assert.Equal(suite.T(), "debug", suite.Conf.Log.Zap.Level)

	// Db
	assert.Equal(suite.T(), "db_apiserver", suite.Conf.Db.DbName)
	assert.Equal(suite.T(), "127.0.0.1", suite.Conf.Db.Host)
	assert.Equal(suite.T(), "3306", suite.Conf.Db.Port)
	assert.Equal(suite.T(), "root", suite.Conf.Db.Username)
	assert.Equal(suite.T(), "root", suite.Conf.Db.Password)
	assert.Equal(suite.T(), true, suite.Conf.Db.LogMode)

	// Mail
	assert.Equal(suite.T(), false, suite.Conf.Mail.Enable)
	assert.Equal(suite.T(), "smtp.exmail.qq.com", suite.Conf.Mail.Smtp)
	assert.Equal(suite.T(), 465, suite.Conf.Mail.Port)
	assert.Equal(suite.T(), "moocss@163.com", suite.Conf.Mail.Username)
	assert.Equal(suite.T(), "", suite.Conf.Mail.Password)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
