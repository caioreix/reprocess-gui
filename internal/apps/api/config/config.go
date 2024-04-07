package config

import (
	"time"

	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
)

// Config represents the configuration structure containing all the project settings.
type Config struct {
	Environment string `default:"dev"`
	Server      Server
	Mongo       Mongo
	Log         Log
}

// Server represents the server configuration.
type Server struct {
	Host              string        `default:"localhost"`
	Port              int           `default:"8080"`
	ReadHeaderTimeout time.Duration `default:"10s"`
}

// Mongo represents the MongoDB configuration.
type Mongo struct {
	URI               string        `required:"true"`
	ConnectionTimeout time.Duration `default:"10s"`
	PingTimeout       time.Duration `default:"2s"`

	TableDatabase   string `required:"true"`
	TableCollection string `required:"true"`

	RowDatabase   string `required:"true"`
	RowCollection string `required:"true"`

	ConsumerDatabase   string `required:"true"`
	ConsumerCollection string `required:"true"`
}

// Log represents the logging configuration.
type Log struct {
	Level string `default:"debug"`
}

// New creates a new Config instance by loading configuration settings from the specified path.
// It loads settings from environment variables and configuration files.
func New(path string) (*Config, error) {
	if path != "" {
		err := godotenv.Load(path)
		if err != nil {
			return nil, err
		}
	}

	config := &Config{}

	settings := &configor.Config{
		Silent:    true,
		ENVPrefix: "API",
	}

	err := configor.New(settings).Load(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
