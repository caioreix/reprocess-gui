package config

import (
	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
)

// Config represents the configuration structure containing all the project settings.
type Config struct {
	Environment string `default:"dev"`
	Log         Log
}

// Log represents the logging configuration.
type Log struct {
	Level string `default:"debug"`
}

// New creates a new Config instance by loading configuration settings from the specified path.
// It loads settings from environment variables and configuration files.
func New(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	settings := &configor.Config{
		Silent:    true,
		ENVPrefix: "WORKER",
	}

	err = configor.New(settings).Load(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
