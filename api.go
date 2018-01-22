package main

import (
	"fmt"
	"github.com/Yapo/logger"
	"github.schibsted.io/Yapo/goms/conf"
	"github.schibsted.io/Yapo/goms/core"
	"github.schibsted.io/Yapo/goms/interfaces"
	"github.schibsted.io/Yapo/goms/repository"
	"github.schibsted.io/Yapo/goms/usecases"
	"gopkg.in/facebookgo/inject.v0"
	"gopkg.in/facebookgo/pidfile.v0"
	"net/http"
	"os"
)

func main() {

	confPath := "conf/conf.json"
	fmt.Printf("Loading config from %s\n", confPath)
	setup, err := conf.Load(confPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conf.Set(setup)

	fmt.Printf("Setting up logger\n")
	loggerConf := logger.LogConfig{
		Syslog: logger.SyslogConfig{
			Enabled:  setup.Logger.SyslogEnabled,
			Identity: setup.Logger.SyslogIdentity,
		},
		Stdlog: logger.StdlogConfig{
			Enabled: setup.Logger.StdlogEnabled,
		},
	}
	if err = logger.Init(loggerConf); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	logger.SetLogLevel(setup.Logger.LogLevel)
	fmt.Printf("LogLevel: %d\n", setup.Logger.LogLevel)

	logger.Info("Saving PID file to %s", setup.Runtime.PidFile)

	pidfile.SetPidfilePath(setup.Runtime.PidFile)
	if err = pidfile.Write(); err != nil {
		logger.Crit("Could not save pid file: %s", err)
		os.Exit(1)
	}

	logger.Info("Setting up Dependency Injection")

	/*
		Setup all the *injectable* resources below.
		Reference: https://godoc.org/github.com/facebookgo/inject
	*/

	var healthHandler interfaces.HealthHandler
	var injectHandler interfaces.InjectHandler

	err = inject.Populate(
		&injectHandler,
		&usecases.ModularCalculator{},
		&repository.ModuloAdder{M: 7},
	)

	if err != nil {
		logger.Crit("%s\n", err)
		os.Exit(1)
	}

	var routes = core.Routes{
		{
			//this is the base path, all routes will start with this
			Prefix: "/api/v{version:[1-9][0-9]*}",
			Groups: []core.Route{
				{
					Name:        "Check service health",
					Method:      "GET",
					Pattern:     "/healthcheck",
					HandlerFunc: healthHandler.Run,
				},
				{
					Name:        "Demonstrate dependency injection with simple math operations",
					Method:      "GET",
					Pattern:     "/inject",
					HandlerFunc: injectHandler.Run,
				},
			},
		},
	}

	logger.Info("Starting request serving")
	logger.Crit("%s\n", http.ListenAndServe(fmt.Sprintf("%s:%d", setup.Runtime.Host, setup.Runtime.Port), core.NewRouter(routes)))
	logger.CloseSyslog() // nolint
}
