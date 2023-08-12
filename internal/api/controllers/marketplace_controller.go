package controllers

import (
	"fmt"
	_ "meteor/docs"
	"meteor/internal/models"
	"strconv"

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

	pickupSity, err := c.service.GetCityByName(bodyy.Pickup.CityName)
	if err != nil {
		c.logger.Error("GetCityByName", zap.Error(err))
		errByte, _ := c.json.Marshal(ErrorResponse{err.Error()})
		c.response(ctx, fasthttp.StatusBadRequest, errByte)
		return
	}
	deliverySity, err := c.service.GetCityByName(bodyy.Delivery.CityName)
	if err != nil {
		c.logger.Error("GetCityByName", zap.Error(err))
		errByte, _ := c.json.Marshal(ErrorResponse{err.Error()})
		c.response(ctx, fasthttp.StatusBadRequest, errByte)
		return
	}
	bodyy.Pickup = pickupSity
	bodyy.Delivery = deliverySity

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

// GetCitiesCashPay
//
//	@Tags			Marketplace
//	@Summary		список городов
//	@Description	список городов
//	@Produce		json
//	@Param			limit query string true "limit городов"
//	@Success		200	{object}	models.ResponseGetCitiesCashPay
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/get-cities-cash-pay [get]
func (c *MarketplaceController) GetCitiesCashPay(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
	limitStr := string(ctx.QueryArgs().Peek("limit"))
	if limitStr == "" {
		limitStr = "0"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		err = fmt.Errorf("strconv.Atoi err: %v", err)
		c.logger.Error("GetCitiesCashPay", zap.Error(err))
		errByte, _ := c.json.Marshal(ErrorResponse{err.Error()})
		c.response(ctx, fasthttp.StatusBadRequest, errByte)
		return
	}
	if limit < 0 {
		err = fmt.Errorf("limit < 0")
		c.logger.Error("GetCitiesCashPay", zap.Error(err))
		errByte, _ := c.json.Marshal(ErrorResponse{err.Error()})
		c.response(ctx, fasthttp.StatusBadRequest, errByte)
		return
	}
	data, err := c.service.GetCitiesCashPay(limit)
	if err != nil {
		err = fmt.Errorf("get order statuses err: %v", err)
		c.logger.Error("GetCitiesCashPay", zap.Error(err))
		errByte, _ := c.json.Marshal(ErrorResponse{err.Error()})
		c.response(ctx, fasthttp.StatusInternalServerError, errByte)
		return
	}

	c.response(ctx, fasthttp.StatusOK, data)
}
