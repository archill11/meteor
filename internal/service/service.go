package service

import (
	"meteor/internal/provider"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type Config struct {
	Proxy DpdCfg `envconfig:"PROXY"`
}

type Service struct {
	cfg      Config
	provider *provider.Provider
	Json     jsoniter.API
	Logger   *zap.Logger
}

func New(cfg Config, provider *provider.Provider, logger *zap.Logger, json jsoniter.API) *Service {
	return &Service{
		cfg:      cfg,
		provider: provider,
		Json:     json,
		Logger:   logger,
	}
}
