package main

import (
	"fmt"
	"github.com/Yapo/logger"
	"github.schibsted.io/Yapo/goms/conf"
	"gopkg.in/facebookgo/inject.v0"
	"gopkg.in/facebookgo/pidfile.v0"
	"net/http"
	"os"
)

var setup *service.Config

func main() {

	confPath := "conf/conf.json"
	fmt.Printf("Loading config from %s\n", confPath)
	setup, err := conf.LoadConf(confPath)
	if err != nil {
		fmt.Printf(err)
		os.Exit(1)
	}
	conf.Set(setup)

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

	logger.Info("Saving PID file to %s", setup.Runtime.PidFile)

	pidfile.SetPidfilePath(setup.Runtime.PidFile)
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
	logger.Crit("%s\n", http.ListenAndServe(fmt.Sprintf("%s:%d", setup.Runtime.Host, setup.Runtime.Port), service.NewRouter()))
	logger.CloseSyslog()
}
