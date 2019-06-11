package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.schibsted.io/Yapo/goms/pkg/infrastructure"
	"github.schibsted.io/Yapo/goms/pkg/interfaces/handlers"

	// CLONE REMOVE START
	"github.schibsted.io/Yapo/goms/pkg/interfaces/loggers"
	"github.schibsted.io/Yapo/goms/pkg/interfaces/repository"
	"github.schibsted.io/Yapo/goms/pkg/usecases"
	// CLONE REMOVE END
)

func main() {
	var shutdownSequence = infrastructure.NewShutdownSequence()
	var conf infrastructure.Config
	shutdownSequence.Listen()
	infrastructure.LoadFromEnv(&conf)
	if jconf, err := json.MarshalIndent(conf, "", "    "); err == nil {
		fmt.Printf("Config: \n%s\n", jconf)
	} else {
		fmt.Printf("Config: \n%+v\n", conf)
	}

	fmt.Printf("Setting up Prometheus\n")
	prometheus := infrastructure.MakePrometheusExporter(
		conf.PrometheusConf.Port,
		conf.PrometheusConf.Enabled,
	)

	fmt.Printf("Setting up logger\n")
	logger, err := infrastructure.MakeYapoLogger(&conf.LoggerConf,
		prometheus.NewEventsCollector(
			"goms_service_events_total",
			"events tracker counter for goms service",
		),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	shutdownSequence.Push(prometheus)

	logger.Info("Initializing resources")

	// HealthHandler
	var healthHandler handlers.HealthHandler

	// CLONE REMOVE START
	// FibonacciHandler
	fibonacciLogger := loggers.MakeFibonacciLogger(logger)
	fibonacciRepository := repository.NewMapFibonacciRepository()
	fibonacciInteractor := usecases.FibonacciInteractor{
		Logger:     fibonacciLogger,
		Repository: fibonacciRepository,
	}
	fibonacciHandler := handlers.FibonacciHandler{
		Interactor: &fibonacciInteractor,
	}

	// To handle http connections you can use an httpHandler or
	// httpCircuitBreakerHandler which retries requests with it's client
	// until it returns a valid answer and then continues normal execution
	// OPTION: classic HTTP
	HTTPHandler := infrastructure.NewHTTPHandler(logger)
	getHealthLogger := loggers.MakeGomsRepoLogger(logger)
	getHealthInteractor := usecases.GetHealthcheckInteractor{
		GomsRepository: repository.NewGomsRepository(
			HTTPHandler,
			conf.GomsClientConf.TimeOut,
			conf.GomsClientConf.GetHealthcheckPath),
		Logger: getHealthLogger,
	}
	getHealthHandler := handlers.GetHealthcheckHandler{
		GetHealthcheckInteractor: &getHealthInteractor,
	}

	// OPTION: HTTP with Circuit Breaker
	circuitBreaker := infrastructure.NewCircuitBreaker(
		conf.CircuitBreakerConf.Name,
		conf.CircuitBreakerConf.ConsecutiveFailure,
		conf.CircuitBreakerConf.FailureRatio,
		conf.CircuitBreakerConf.Timeout,
		conf.CircuitBreakerConf.Interval,
		logger,
	)
	HTTPCBHandler := infrastructure.NewHTTPCircuitBreakerHandler(circuitBreaker, logger, HTTPHandler)
	getHealthCBHandler := handlers.GetHealthcheckHandler{
		GetHealthcheckInteractor: &usecases.GetHealthcheckInteractor{
			GomsRepository: repository.NewGomsRepository(
				HTTPCBHandler,
				conf.GomsClientConf.TimeOut,
				conf.GomsClientConf.GetHealthcheckPath),
			Logger: getHealthLogger,
		},
	}
	// CLONE REMOVE END

	// Setting up router
	maker := infrastructure.RouterMaker{
		Logger:        logger,
		WrapperFuncs:  []infrastructure.WrapperFunc{prometheus.TrackHandlerFunc},
		WithProfiling: conf.ServiceConf.Profiling,
		Routes: infrastructure.Routes{
			{
				// This is the base path, all routes will start with this prefix
				Prefix: "/api/v{version:[1-9][0-9]*}",
				Groups: []infrastructure.Route{
					{
						Name:    "Check service health",
						Method:  "GET",
						Pattern: "/healthcheck",
						Handler: &healthHandler,
					},
					// CLONE REMOVE START
					{
						Name:    "Retrieve the Nth Fibonacci with Clean Architecture",
						Method:  "GET",
						Pattern: "/fibonacci",
						Handler: &fibonacciHandler,
					},
					{
						Name:    "Retrieve healthcheck by doing a client request to itself",
						Method:  "GET",
						Pattern: "/http/healthcheck",
						Handler: &getHealthHandler,
					},
					{
						Name:    "Retrieve healthcheck by doing a client request to itself using Circuit Breaker",
						Method:  "GET",
						Pattern: "/httpcb/healthcheck",
						Handler: &getHealthCBHandler,
					},
					// CLONE REMOVE END
				},
			},
		},
	}

	router := maker.NewRouter()

	server := infrastructure.NewHTTPServer(
		fmt.Sprintf("%s:%d", conf.Runtime.Host, conf.Runtime.Port),
		router,
		logger,
	)
	shutdownSequence.Push(server)
	logger.Info("Starting request serving")
	go server.ListenAndServe()
	shutdownSequence.Wait()
	logger.Info("Server exited normally")

}
