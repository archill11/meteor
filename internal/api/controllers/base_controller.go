package controllers

import (
	"meteor/internal/service"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type BaseController struct {
	service *service.Service
	logger  *zap.Logger
	json    jsoniter.API
}

func NewBaseController(service *service.Service) *BaseController {
	return &BaseController{
		service,
		service.Logger,
		service.Json,
	}
}

func (c *BaseController) response(ctx *fasthttp.RequestCtx, code int, data []byte) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(code)

	_, err := ctx.Write(data)
	if err != nil {
		c.logger.Error("Failed to write data into response body", zap.Error(err), zap.String("action", "SEND_RESPONSE"))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
}

type ErrorResponse struct {
	Error any `json:"error"`
}
