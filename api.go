package main

import (
	"fmt"
	"github.com/Yapo/logger"
	"github.schibsted.io/Yapo/goms/service"
	"gopkg.in/facebookgo/inject.v0"
	"gopkg.in/facebookgo/pidfile.v0"
	"net/http"
	"os"
)

var setup *service.Config

func main() {

	service.LoadConf("conf/conf.json")
	setup = service.GetConfig()

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

	logger.Info("Saving PID file to %s", setup.Service.PidFile)

	pidfile.SetPidfilePath(setup.Service.PidFile)
	if err := pidfile.Write(); err != nil {
		logger.Crit("Could not save pid file: %s", err)
		os.Exit(1)
	}

	logger.Info("Setting up Dependency Injection")

	/*
		Setup all the *injectable* resources below.
		Reference: https://godoc.org/github.com/facebookgo/inject
	*/
	err := service.SetupInject(
		&inject.Object{Value: &service.Resource{Name: "A Resource", Usage: "Being injected"}},
	)
	if err != nil {
		logger.Crit("%s\n", err)
		os.Exit(1)
	}

	logger.Info("Starting request serving")
	logger.Crit("%s\n", http.ListenAndServe(fmt.Sprintf("%s:%d", setup.Service.Host, setup.Service.Port), service.NewRouter()))
	logger.CloseSyslog()
}
