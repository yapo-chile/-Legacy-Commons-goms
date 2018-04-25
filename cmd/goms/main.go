package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Yapo/logger"
	"github.schibsted.io/Yapo/goms/pkg/core"
	"github.schibsted.io/Yapo/goms/pkg/interfaces/handlers"
	"github.schibsted.io/Yapo/goms/pkg/interfaces/repository"
	"github.schibsted.io/Yapo/goms/pkg/usecases"
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

	// HealthHandler
	var healthHandler handlers.HealthHandler

	// FibonacciHandler
	fibonacciRepository := repository.NewMapFibonacciRepository()
	fibonacciInteractor := usecases.FibonacciInteractor{
		Repository: fibonacciRepository,
	}
	fibonacciHandler := handlers.FibonacciHandler{
		Interactor: &fibonacciInteractor,
	}

	var routes = core.Routes{
		{
			//this is the base path, all routes will start with this
			Prefix: "/api/v{version:[1-9][0-9]*}",
			Groups: []core.Route{
				{
					Name:    "Check service health",
					Method:  "GET",
					Pattern: "/healthcheck",
					Handler: &healthHandler,
				},
				{
					Name:    "Retrieve the Nth Fibonacci with Clean Architecture",
					Method:  "GET",
					Pattern: "/fibonacci",
					Handler: &fibonacciHandler,
				},
			},
		},
	}

	logger.Info("Starting request serving")
	logger.Crit("%s\n", http.ListenAndServe(fmt.Sprintf("%s:%d", conf.AppConf.Host, conf.AppConf.Port), core.NewRouter(routes)))
	logger.CloseSyslog() // nolint
}
