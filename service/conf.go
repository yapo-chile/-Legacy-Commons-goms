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

type ServiceConfig struct {
	Host    string
	Port    int
	PidFile string
}

/*
	LogLevel definition:
	a0 - Debug
	1 - Info
	2 - Warning
	3 - Error
	4 - Critic
*/
type LoggerConfig struct {
	SyslogEnabled  bool
	SyslogIdentity string
	StdlogEnabled  bool
	LogLevel       int
}

type Config struct {
	Service ServiceConfig
	Logger  LoggerConfig
}

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

func SetConf(c *Config) {
	config = c
}

func GetConfig() *Config {
	return config
}
