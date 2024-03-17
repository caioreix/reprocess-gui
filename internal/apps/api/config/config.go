package config

import (
	"time"

	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string `default:"dev"`
	Server      Server
	Mongo       Mongo
	Log         Log
}

type Server struct {
	Host string `default:"localhost"`
	Port int    `default:"8080"`
}

type Mongo struct {
	URI               string        `required:"true"`
	ConnectionTimeout time.Duration `default:"10s"`
	PingTimeout       time.Duration `default:"2s"`

	TableDatabase   string `required:"true"`
	TableCollection string `required:"true"`
}

type Log struct {
	Level string `default:"debug"`
}

func New(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	settings := &configor.Config{
		Silent:    true,
		ENVPrefix: "API",
	}

	err = configor.New(settings).Load(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
