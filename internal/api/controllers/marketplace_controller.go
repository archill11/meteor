package controllers

import (
	"fmt"
	_ "meteor/docs"
	"meteor/internal/models"

	"go.uber.org/zap"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

type MarketplaceController struct {
	*BaseController
}

func NewMarketplaceControllerFromBase(base *BaseController) *MarketplaceController {
	return &MarketplaceController{
		base,
	}
}

// GetServiceCost
//
//	@Tags			Marketplace
//	@Summary		рассчитать стоимость доставки по параметрам посылок
//	@Description	рассчитать стоимость доставки по параметрам посылок
//	@Produce		json
//	@Accept			json
//	@Param			ServiceData	body		models.RequestGetServiceCost	true	"данные заказа"
//	@Success		200			{object}	models.ResponseGetServiceCost
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/v1/get-service-cost [post]
func (c *MarketplaceController) GetServiceCost(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {

	var bodyy models.RequestGetServiceCost
	if err := c.json.Unmarshal(ctx.Request.Body(), &bodyy); err != nil {
		err := fmt.Errorf("can't decode body: %v", err)
		c.logger.Error("GetServiceCost", zap.Error(err))
		errByte, _ := c.json.Marshal(ErrorResponse{err.Error()})
		c.response(ctx, fasthttp.StatusBadRequest, errByte)
		return
	}

	data, err := c.service.GetServiceCost(ctx, bodyy)
	if err != nil {
		err = fmt.Errorf("get order statuses err: %v", err)
		c.logger.Error("GetServiceCost", zap.Error(err))
		errByte, _ := c.json.Marshal(ErrorResponse{err.Error()})
		c.response(ctx, fasthttp.StatusInternalServerError, errByte)
		return
	}

	c.response(ctx, fasthttp.StatusOK, data)
}
