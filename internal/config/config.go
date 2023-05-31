package config

import (
	"fmt"
	"meteor/internal/api/server"
	"meteor/internal/provider"
	"meteor/internal/service"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   server.ServerConfig `envconfig:"SERVER"`
	Service  service.Config      `envconfig:"SERVICE"`
	Provider provider.Config     `envconfig:"PROVIDER"`
}

func New() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}
	
	var cfg = new(Config)

	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
