package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.schibsted.io/Yapo/goms/pkg/infrastructure"
	"github.schibsted.io/Yapo/goms/pkg/interfaces/handlers"
	"github.schibsted.io/Yapo/goms/pkg/interfaces/repository"
	"github.schibsted.io/Yapo/goms/pkg/usecases"
)

func main() {
	var conf infrastructure.Config
	infrastructure.LoadFromEnv(&conf)
	jconf, _ := json.MarshalIndent(conf, "", "    ")
	fmt.Printf("Config:\n%s\n", jconf)

	fmt.Printf("Setting up logger\n")
	logger, err := infrastructure.MakeYapoLogger(&conf.LoggerConf)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	logger.Info("Initializing resources")

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

	var routes = infrastructure.Routes{
		{
			//this is the base path, all routes will start with this
			Prefix: "/api/v{version:[1-9][0-9]*}",
			Groups: []infrastructure.Route{
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
	logger.Crit("%s\n", http.ListenAndServe(conf.ServiceConf.Host, infrastructure.NewRouter(routes)))
}
