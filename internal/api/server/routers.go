package server

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

type Routers interface {
	HealthCheck(*fasthttp.RequestCtx, fasthttprouter.Params)
	HandleDocs(*fasthttp.RequestCtx, fasthttprouter.Params)

	GetServiceCost(*fasthttp.RequestCtx, fasthttprouter.Params)
}

func NewRouters(routers Routers) func(ctx *fasthttp.RequestCtx) {
	router := fasthttprouter.New()
	Handler := router.Handler

	router.GET("/health", routers.HealthCheck)
	router.GET("/swagger/*filepath", routers.HandleDocs)

	router.POST("/api/v1/get-service-cost", routers.GetServiceCost)

	return Handler
}
