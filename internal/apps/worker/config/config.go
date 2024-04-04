package config

import (
	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
)

type Config struct {
	Environment string `default:"dev"`
	Log         Log
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
