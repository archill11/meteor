package controllers

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/valyala/fasthttprouter"

	_ "meteor/docs"
)

type MainController struct {
	*BaseController
}

func NewMainControllerFromBase(base *BaseController) *MainController {
	return &MainController{
		base,
	}
}

// HealthCheck
//
//	@Tags			Health
//	@Summary		проверка работоспособности
//	@Description	Проверка работоспособности
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		500	{object}	ErrorResponse
//	@Router			/health [get]
func (r *MainController) HealthCheck(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (r *MainController) HandleDocs(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	fasthttpadaptor.NewFastHTTPHandler(httpSwagger.Handler())(ctx)
}

func (r *MainController) HandleMetrics(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())(ctx)
}
