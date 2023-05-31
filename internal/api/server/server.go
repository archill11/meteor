package server

import (
	"time"

	"github.com/valyala/fasthttp"
)

type ServerConfig struct {
	ListenAddress string        `default:":9090" envconfig:"PORT"`
	WriteTimeout  time.Duration `default:"1m" envconfig:"WRITE_TIMEOUT"`
	ReadTimeout   time.Duration `default:"1m" envconfig:"READ_TIMEOUT"`
}

func NewServer(cfg ServerConfig, handler func(ctx *fasthttp.RequestCtx)) *fasthttp.Server {
	return &fasthttp.Server{
		Name:         "meteor",
		Handler:      handler,
		WriteTimeout: time.Minute*20,
		ReadTimeout:  time.Minute*20,
	}
}
