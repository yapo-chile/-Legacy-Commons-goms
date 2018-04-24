package core

import (
	"fmt"
	"os"
	"strconv"
)

//LoggerConf configuration for the logger struct used for logging
// LogLevel definition:
// 0 - Debug
// 1 - Info
// 2 - Warning
// 3 - Error
// 4 - Critic
type LoggerConf struct {
	SyslogIdentity string
	LogLevel       int
	SyslogEnabled  bool
	StdlogEnabled  bool
}

//ServiceConf struct representing the config of a service (redis, postgresql, etc)
type ServiceConf struct {
	Host    string
	Port    int
	Version string
}

//Config struct that represents the config of the microservice
type Config struct {
	AppConf    ServiceConf
	LoggerConf LoggerConf
}

//Load load the microservice conf using the enviroment variables
func Load() (*Config, error) {
	conf := new(Config)
	port, err := strconv.Atoi(os.Getenv("LISTEN_PORT"))
	if err != nil {
		return nil, fmt.Errorf("Error loading config %+v", err)
	}
	conf.AppConf = ServiceConf{
		Host:    os.Getenv("SERVERNAME"),
		Port:    port,
		Version: os.Getenv("VERSION"),
	}
	syslogEnabled, err := strconv.ParseBool(os.Getenv("SYSLOG_ENABLED"))
	if err != nil {
		return nil, fmt.Errorf("Error loading config %+v", err)
	}
	stdLogEnabled, err := strconv.ParseBool(os.Getenv("STDLOG_ENABLED"))
	if err != nil {
		return nil, fmt.Errorf("Error loading config %+v", err)
	}
	logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return nil, fmt.Errorf("Error loading config %+v", err)
	}

	conf.LoggerConf = LoggerConf{
		SyslogEnabled:  syslogEnabled,
		SyslogIdentity: os.Getenv("SYSLOG_IDENTITY"),
		StdlogEnabled:  stdLogEnabled,
		LogLevel:       logLevel,
	}

	return conf, nil
}
