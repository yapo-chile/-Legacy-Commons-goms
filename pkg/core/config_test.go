package core

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestConfigLoad(t *testing.T) {
	configVariables := []string{
		"SERVERNAME",
		"LISTEN_PORT",
		"VERSION",
		"SYSLOG_ENABLED",
		"SYSLOG_IDENTITY",
		"STDLOG_ENABLED",
		"LOG_LEVEL",
	}
	storedValues := make([]string, len(configVariables))
	testConfigVariables := []string{
		"localhost",
		"8081",
		"test_version",
		"true",
		"test_log",
		"true",
		"1",
	}
	//store the environ variables and set the desired values
	for index, variable := range configVariables {
		storedValues[index] = os.Getenv(variable)
		os.Setenv(variable, testConfigVariables[index])
	}
	//Load the test enviroment
	c, err := Load()
	for index, variable := range configVariables {
		os.Setenv(variable, storedValues[index])
	}
	//Validate the load
	assert.Equal(t, err, nil)
	assert.Equal(t, fmt.Sprintf("%+v", c), `&{AppConf:{Host:localhost Port:8081 Version:test_version} LoggerConf:{SyslogEnabled:true SyslogIdentity:test_log StdlogEnabled:true LogLevel:1}}`)
}

func TestConfigLoadNoPort(t *testing.T) {
	configVariables := []string{
		"SERVERNAME",
		"LISTEN_PORT",
		"VERSION",
		"SYSLOG_ENABLED",
		"SYSLOG_IDENTITY",
		"STDLOG_ENABLED",
		"LOG_LEVEL",
	}
	storedValues := make([]string, len(configVariables))
	testConfigVariables := []string{
		"localhost",
		"",
		"test_version",
		"true",
		"test_log",
		"true",
		"1",
	}
	//store the environ variables and set the desired values
	for index, variable := range configVariables {
		storedValues[index] = os.Getenv(variable)
		os.Setenv(variable, testConfigVariables[index])
	}
	//Load the test enviroment
	_, err := Load()
	for index, variable := range configVariables {
		os.Setenv(variable, storedValues[index])
	}
	//Validate the load
	assert.Equal(t, fmt.Sprintf("%+v", err), "Error loading config strconv.Atoi: parsing \"\": invalid syntax")
}

func TestConfigLoadNoSyslog(t *testing.T) {
	configVariables := []string{
		"SERVERNAME",
		"LISTEN_PORT",
		"VERSION",
		"SYSLOG_ENABLED",
		"SYSLOG_IDENTITY",
		"STDLOG_ENABLED",
		"LOG_LEVEL",
	}
	storedValues := make([]string, len(configVariables))
	testConfigVariables := []string{
		"localhost",
		"8081",
		"test_version",
		"",
		"test_log",
		"true",
		"1",
	}
	//store the environ variables and set the desired values
	for index, variable := range configVariables {
		storedValues[index] = os.Getenv(variable)
		os.Setenv(variable, testConfigVariables[index])
	}
	//Load the test enviroment
	_, err := Load()
	for index, variable := range configVariables {
		os.Setenv(variable, storedValues[index])
	}
	//Validate the load
	assert.Equal(t, fmt.Sprintf("%+v", err), "Error loading config strconv.ParseBool: parsing \"\": invalid syntax")
}

func TestConfigLoadNoStdlog(t *testing.T) {
	configVariables := []string{
		"SERVERNAME",
		"LISTEN_PORT",
		"VERSION",
		"SYSLOG_ENABLED",
		"SYSLOG_IDENTITY",
		"STDLOG_ENABLED",
		"LOG_LEVEL",
	}
	storedValues := make([]string, len(configVariables))
	testConfigVariables := []string{
		"localhost",
		"8081",
		"test_version",
		"true",
		"test_log",
		"",
		"1",
	}
	//store the environ variables and set the desired values
	for index, variable := range configVariables {
		storedValues[index] = os.Getenv(variable)
		os.Setenv(variable, testConfigVariables[index])
	}
	//Load the test enviroment
	_, err := Load()
	for index, variable := range configVariables {
		os.Setenv(variable, storedValues[index])
	}
	//Validate the load
	assert.Equal(t, fmt.Sprintf("%+v", err), "Error loading config strconv.ParseBool: parsing \"\": invalid syntax")
}

func TestConfigLoadNoLogLevel(t *testing.T) {
	configVariables := []string{
		"SERVERNAME",
		"LISTEN_PORT",
		"VERSION",
		"SYSLOG_ENABLED",
		"SYSLOG_IDENTITY",
		"STDLOG_ENABLED",
		"LOG_LEVEL",
	}
	storedValues := make([]string, len(configVariables))
	testConfigVariables := []string{
		"localhost",
		"8081",
		"test_version",
		"true",
		"test_log",
		"true",
		"",
	}
	//store the environ variables and set the desired values
	for index, variable := range configVariables {
		storedValues[index] = os.Getenv(variable)
		os.Setenv(variable, testConfigVariables[index])
	}
	//Load the test enviroment
	_, err := Load()
	for index, variable := range configVariables {
		os.Setenv(variable, storedValues[index])
	}
	//Validate the load
	assert.Equal(t, fmt.Sprintf("%+v", err), "Error loading config strconv.Atoi: parsing \"\": invalid syntax")
}