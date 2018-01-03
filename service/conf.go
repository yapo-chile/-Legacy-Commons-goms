package service

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

// LoadConf retrieves the configuration from the file specified by path.
// Sets the global conf to this value.
func LoadConf(path string) {
	fmt.Printf("Loading config from %s\n", path)

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error loading file %s, error: %s\n", path, err)
	}
	defer file.Close()

	fmt.Printf("Decoding conf file\n")
	decoder := json.NewDecoder(file)
	c := &Config{}

	if err := decoder.Decode(&c); err != nil {
		fmt.Printf("Error decoding conf file: %s, error: %s\n", configPath, err)
	}
	SetConf(c)
}

// SetConf updates the global configuration to given value.
func SetConf(c *Config) {
	config = c
}

// GetConfig retrieves the global configuration.
func GetConfig() *Config {
	return config
}
