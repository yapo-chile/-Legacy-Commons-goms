package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	configPath string
	config     *Config
)

// RuntimeConfig allows magic parsing of the Runtime directive on conf.json
type RuntimeConfig struct {
	Host    string
	Port    int
	PidFile string
}

// LoggerConfig allows magic parsing of the Logger directive on conf.json
// LogLevel definition:
// 0 - Debug
// 1 - Info
// 2 - Warning
// 3 - Error
// 4 - Critic
type LoggerConfig struct {
	SyslogEnabled  bool
	SyslogIdentity string
	StdlogEnabled  bool
	LogLevel       int
}

// Config type for parsing configuration on conf.json
type Config struct {
	Runtime RuntimeConfig
	Logger  LoggerConfig
}

// Load retrieves the configuration from the file specified by path.
// Sets the global conf to this value.
func Load(path string) (*Config, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("conf: Error loading file %s, error: %s\n", path, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	c := &Config{}

	if err := decoder.Decode(c); err != nil {
		return nil, fmt.Errorf("conf: Error decoding conf file: %s, error: %s\n", configPath, err)
	}
	return c, nil
}

// Set updates the global configuration to given value.
func Set(c *Config) {
	config = c
}

// Get retrieves the global configuration.
func Get() *Config {
	return config
}
