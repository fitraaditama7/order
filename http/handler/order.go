package handler

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"order-backend/constant"
	"order-backend/order"
	"order-backend/pkg/customerror"
	"order-backend/pkg/logger"
	"order-backend/pkg/utils"
	"strconv"
	"time"
)

type OrderService interface {
	GetListOrder(ctx context.Context, param order.Param) (*order.OrderResponse, error)
}

type OrderHandler struct {
	orderService OrderService
}

func NewOrderHandler(orderService OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (o *OrderHandler) GetList(c *fiber.Ctx) error {
	log := logger.Log()
	var err error
	var dateStart *time.Time
	var dateEnd *time.Time
	limit := constant.DEFAULT_LIMIT
	page := constant.DEFAULT_PAGE
	keyword := c.Query("keyword")

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			log.Error(err.Error())
			return ErrorResponse(c, fiber.StatusBadRequest, customerror.Error(customerror.ErrCodeBadRequest, customerror.ErrMessageBadRequest))
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			log.Error(err.Error())
			return ErrorResponse(c, fiber.StatusBadRequest, customerror.Error(customerror.ErrCodeBadRequest, customerror.ErrMessageBadRequest))
		}
	}

	if c.Query("date_start") != "" && c.Query("date_end") != "" {
		dateStart, err = utils.FormatStartDate(c.Query("date_start"))
		if err != nil {
			log.Error(err.Error())
			return ErrorResponse(c, fiber.StatusBadRequest, customerror.Error(customerror.ErrCodeDateFormat, customerror.ErrMessageDateFormat))
		}

		dateEnd, err = utils.FormatEndDate(c.Query("date_end"))
		if err != nil {
			log.Error(err.Error())
			return ErrorResponse(c, fiber.StatusBadRequest, customerror.Error(customerror.ErrCodeDateFormat, customerror.ErrMessageDateFormat))
		}
	}

	result, err := o.orderService.GetListOrder(c.Context(), order.Param{
		Limit:     limit,
		Page:      page,
		Keyword:   keyword,
		DateStart: dateStart,
		DateEnd:   dateEnd,
	})
	if err != nil {
		log.Error(err.Error())
		return ErrorResponse(c, fiber.StatusInternalServerError, customerror.Error(customerror.ErrCodeGeneric, err.Error()))
	}

	return SuccessResponse(c, result)
}
