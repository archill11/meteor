package service

import (
	"context"
	"fmt"
	"meteor/internal/models"

	"github.com/valyala/fasthttp"
)

type DpdCfg struct {
	Host string `default:"" envconfig:"HOST"`
}

func (s *Service) GetServiceCost(ctx *fasthttp.RequestCtx, body models.RequestGetServiceCost) ([]byte, error) {
	requestUri := string(ctx.Request.RequestURI())
	fmt.Println()
	newReqUri := fmt.Sprintf("%s%s%s", s.cfg.Proxy.Host, "", requestUri)
	result, _, err := s.provider.SendRequest(context.Background(), ctx, newReqUri, nil)
	if err != nil {
		return nil, fmt.Errorf("GetServiceCost http.Post: %v", err)
	}
	return result, nil
}
