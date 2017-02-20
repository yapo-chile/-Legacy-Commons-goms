package main

import (
	"fmt"
	"github.com/Yapo/logger"
	"net/http"
)

var setup *Config

func main() {

	Load("conf/conf.json")
	setup = GetConfig()

	loggerConf := logger.LogConfig{
		Syslog: logger.SyslogConfig{
			Enabled:  setup.Logger.SyslogEnabled,
			Identity: setup.Logger.SyslogIdentity,
		},
		Stdlog: logger.StdlogConfig{
			Enabled: setup.Logger.StdlogEnabled,
		},
	}
	logger.Init(loggerConf)
	logger.SetLogLevel(setup.Logger.LogLevel)
	fmt.Printf("LogLevel: %d\n", setup.Logger.LogLevel)

	logger.Crit("%s\n", http.ListenAndServe(fmt.Sprintf("%s:%d", setup.Service.Host, setup.Service.Port), NewRouter()))
	logger.CloseSyslog()
}
