package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	jsoniter "github.com/json-iterator/go"
)

type Config struct {
	WriteTimeout time.Duration `default:"1m" envconfig:"WRITE_TIMEOUT"`
	ReadTimeout  time.Duration `default:"1m" envconfig:"READ_TIMEOUT"`
}

type Provider struct {
	p      *fasthttp.Client
	cfg    Config
	json   jsoniter.API
	logger *zap.Logger
}

func New(cfg Config, logger *zap.Logger, json jsoniter.API) *Provider {
	return &Provider{
		p: &fasthttp.Client{
			WriteTimeout: cfg.WriteTimeout,
			ReadTimeout:  cfg.ReadTimeout,
		},
		cfg:    cfg,
		json:   json,
		logger: logger,
	}
}

func (p *Provider) SendRequest(c context.Context, ctx *fasthttp.RequestCtx, url string, fn func(*fasthttp.Request)) ([]byte, int, error) {
	err := c.Err()
	if err != nil {
		return nil, fasthttp.StatusBadRequest, err
	}

	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()
	ctx.Request.CopyTo(req)
	req.SetRequestURI(url)
	if fn != nil {
		fn(req)
	}

	if err := p.p.Do(req, res); err != nil {
		return nil, fasthttp.StatusInternalServerError, fmt.Errorf("send request: %v", err)
	}

	return res.Body(), res.StatusCode(), nil
}
