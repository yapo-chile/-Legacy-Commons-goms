package infrastructure

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// ServiceConf holds configuration for this Service
type ServiceConf struct {
	Host      string `env:"HOST" envDefault:":8080"`
	Profiling bool   `env:"PROFILING" envDefault:"true"`
}

// LoggerConf holds configuration for logging
// LogLevel definition:
//   0 - Debug
//   1 - Info
//   2 - Warning
//   3 - Error
//   4 - Critic
type LoggerConf struct {
	SyslogIdentity string `env:"SYSLOG_IDENTITY"`
	SyslogEnabled  bool   `env:"SYSLOG_ENABLED" envDefault:"false"`
	StdlogEnabled  bool   `env:"STDLOG_ENABLED" envDefault:"true"`
	LogLevel       int    `env:"LOG_LEVEL" envDefault:"0"`
}

// PrometheusConf holds configuration to report to Prometheus
type PrometheusConf struct {
	Port    string `env:"PORT" envDefault:"8877"`
	Enabled bool   `env:"ENABLED" envDefault:"false"`
}

// RuntimeConfig config to start the app
type RuntimeConfig struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port int    `env:"PORT" envDefault:"8080"`
}

// CircuitBreakerConf holds all configurations for circuit breaker
type CircuitBreakerConf struct {
	Name               string  `env:"NAME" envDefault:"HTTP_SEND"`
	ConsecutiveFailure uint32  `env:"CONSECUTIVE_FAILURE" envDefault:"10"`
	FailureRatio       float64 `env:"FAILURE_RATIO" envDefault:"0.5"`
	Timeout            int     `env:"TIMEOUT" envDefault:"30"`
	Interval           int     `env:"INTERVAL" envDefault:"30"`
}

// GomsClientConf holds configuration regarding to our http client (goms itself in this case)
type GomsClientConf struct {
	TimeOut            int    `env:"TIMEOUT" envDefault:"30"`
	GetHealthcheckPath string `env:"HEALTH_PATH" envDefault:"/get/healthcheck"`
}

// EtcdConf configure how to read configuration from remote Etcd service
type EtcdConf struct {
	Host       string `env:"HOST" envDefault:"http://lb:2397"`
	LastUpdate string `env:"LAST_UPDATE" envDefault:"/last_update"`
	Prefix     string `env:"PREFIX" envDefault:"/v2/keys"`
}

// Config holds all configuration for the service
type Config struct {
	ServiceConf        ServiceConf        `env:"SERVICE_"`
	PrometheusConf     PrometheusConf     `env:"PROMETHEUS_"`
	LoggerConf         LoggerConf         `env:"LOGGER_"`
	Runtime            RuntimeConfig      `env:"APP_"`
	CircuitBreakerConf CircuitBreakerConf `env:"CIRCUIT_BREAKER_"`
	GomsClientConf     GomsClientConf     `env:"GOMS_"`
	EtcdConf           EtcdConf           `env:"ETCD_"`
}

// LoadFromEnv loads the config data from the environment variables
func LoadFromEnv(data interface{}) {
	load(reflect.ValueOf(data), "", "")
}

// valueFromEnv lookup the best value for a variable on the environment
func valueFromEnv(envTag, envDefault string) string {
	// Maybe it's a secret and <envTag>_FILE points to a file with the value
	// https://rancher.com/docs/rancher/v1.6/en/cattle/secrets/#docker-hub-images
	if fileName, ok := os.LookupEnv(fmt.Sprintf("%s_FILE", envTag)); ok {
		// filepath.Clean() will clean the input path and remove some unnecessary things
		// like multiple separators doble "." and others
		// if for some reason you are having troubles reaching your file, check the
		// output of the Clean function and test if its what you expect
		// you can find more info here: https://golang.org/pkg/path/filepath/#Clean
		b, err := ioutil.ReadFile(filepath.Clean(fileName))
		if err == nil {
			return string(b)
		}
		fmt.Print(err)
	}
	// The value might be set directly on the environment
	if value, ok := os.LookupEnv(envTag); ok {
		return value
	}
	// Nothing to do, return the default
	return envDefault
}

// load the variable defined in the envTag into Value
func load(conf reflect.Value, envTag, envDefault string) {
	if conf.Kind() == reflect.Ptr {
		reflectedConf := reflect.Indirect(conf)
		// Only attempt to set writeable variables
		if reflectedConf.IsValid() && reflectedConf.CanSet() {
			value := valueFromEnv(envTag, envDefault)
			// Print message if config is missing
			if envTag != "" && value == "" && !strings.HasSuffix(envTag, "_") {
				fmt.Printf("Config for %s missing\n", envTag)
			}
			switch reflectedConf.Kind() {
			case reflect.Struct:
				// Recursively load inner struct fields
				for i := 0; i < reflectedConf.NumField(); i++ {
					if tag, ok := reflectedConf.Type().Field(i).Tag.Lookup("env"); ok {
						def, _ := reflectedConf.Type().Field(i).Tag.Lookup("envDefault")
						load(reflectedConf.Field(i).Addr(), envTag+tag, def)
					}
				}
			// Here for each type we should make a cast of the env variable and then set the value
			case reflect.String:
				reflectedConf.SetString(value)
			case reflect.Int:
				if value, err := strconv.Atoi(value); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case reflect.Float64:
				if value, err := strconv.ParseFloat(value, 64); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case reflect.Uint32:
				if value, err := strconv.ParseUint(value, 10, 32); err == nil {
					reflectedConf.Set(reflect.ValueOf(uint32(value)))
				}
			case reflect.Bool:
				if value, err := strconv.ParseBool(value); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			}
		}
	}
}
