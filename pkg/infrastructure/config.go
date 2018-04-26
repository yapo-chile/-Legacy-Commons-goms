package infrastructure

import (
	"os"
	"reflect"
	"strconv"
)

// ServiceConf holds configuration for this Service
type ServiceConf struct {
	Host string `env:"HOST" envDefault:":8080"`
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

// NewRelicConf holds configuration to report to New Relic
// TODO: You need to set the defaults according to your service
type NewRelicConf struct {
	Key     string `env:"KEY" envDefault:"923864cba2f12410aff39279cddfd339a07f13a3"`
	Appname string `env:"APPNAME" envDefault:"yapo-goms-poya"`
}

//Config holds all configuration for the service
type Config struct {
	ServiceConf  ServiceConf  `env:"SERVICE_"`
	NewRelicConf NewRelicConf `env:"NEWRELIC_"`
	LoggerConf   LoggerConf   `env:"LOGGER_"`
}

// LoadConfig loads the config data from the environment variables
func LoadFromEnv(data interface{}) {
	load(reflect.ValueOf(data), "", "")
}

//load loads the variable defined in the envTag into the value
func load(conf reflect.Value, envTag, envDefault string) {
	if conf.Kind() == reflect.Ptr {
		reflectedConf := reflect.Indirect(conf)
		// we should only keep going if we can set values
		if reflectedConf.IsValid() && reflectedConf.CanSet() {
			value, ok := os.LookupEnv(envTag)
			// if the env variable is not set we just use the envDefault
			if !ok {
				value = envDefault
			}
			switch reflectedConf.Kind() {
			case reflect.Struct:
				for i := 0; i < reflectedConf.NumField(); i++ {
					if tag, ok := reflectedConf.Type().Field(i).Tag.Lookup("env"); ok {
						def, _ := reflectedConf.Type().Field(i).Tag.Lookup("envDefault")
						load(reflectedConf.Field(i).Addr(), envTag+tag, def)
					}
				}
				break
			// Here for each type we should make a cast of the env variable and then set the value
			case reflect.String:
				reflectedConf.SetString(value)
				break
			case reflect.Int:
				value, _ := strconv.Atoi(value)
				reflectedConf.Set(reflect.ValueOf(value))
				break
			case reflect.Bool:
				value, _ := strconv.ParseBool(value)
				reflectedConf.Set(reflect.ValueOf(value))
			}
		}
	}
}
