package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Yapo/logger"
	"github.schibsted.io/Yapo/goms/pkg/core"
	"github.schibsted.io/Yapo/goms/pkg/interfaces"
	"github.schibsted.io/Yapo/goms/pkg/repository"
	"github.schibsted.io/Yapo/goms/pkg/usecases"
	"gopkg.in/facebookgo/inject.v0"
)

func main() {

	fmt.Printf("Loading config")
	conf, err := core.Load()
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(2)
	}
	fmt.Printf("Loaded config %v\n", conf)
	fmt.Printf("Setting up logger\n")
	loggerConf := logger.LogConfig{
		Syslog: logger.SyslogConfig{
			Enabled:  conf.LoggerConf.SyslogEnabled,
			Identity: conf.LoggerConf.SyslogIdentity,
		},
		Stdlog: logger.StdlogConfig{
			Enabled: conf.LoggerConf.StdlogEnabled,
		},
	}
	if err := logger.Init(loggerConf); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	logger.SetLogLevel(conf.LoggerConf.LogLevel)
	fmt.Printf("LogLevel: %d\n", conf.LoggerConf.LogLevel)

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
	logger.Crit("%s\n", http.ListenAndServe(fmt.Sprintf("%s:%d", conf.AppConf.Host, conf.AppConf.Port), core.NewRouter(routes)))
	logger.CloseSyslog() // nolint
}
