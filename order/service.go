package order

import (
	"context"
	"fmt"
	"log"
	"math"
	"order-backend/constant"
	"order-backend/pkg/customerror"
	"order-backend/pkg/logger"
	"order-backend/repository"
	"time"
)

type service struct {
	orderRepository OrderRepository
}

type OrderRepository interface {
	FindList(ctx context.Context, limit, page int, keyword string, dateStart, dateEnd *time.Time) ([]repository.OrderList, error)
	FindTotalOrder(ctx context.Context, keyword string, dateStart, dateEnd *time.Time) (*repository.TotalOrder, error)
}

func NewOrderService(orderRepository OrderRepository) *service {
	return &service{orderRepository: orderRepository}
}

func (s *service) GetListOrder(ctx context.Context, param Param) (*OrderResponse, error) {
	log := logger.Log()

	orders, err := s.orderRepository.FindList(ctx, param.Limit, param.Page, param.Keyword, param.DateStart, param.DateEnd)
	if err != nil {
		log.Error(err.Error())
		return nil, customerror.Error(customerror.ErrCodeGeneric, customerror.ErrMessageGeneric)
	}

	if len(orders) == 0 {
		data := make([]OrderListResponse, 0)
		return &OrderResponse{CurrentPage: 1, TotalPage: 1, Data: data}, nil
	}

	totalOrder, err := s.orderRepository.FindTotalOrder(ctx, param.Keyword, param.DateStart, param.DateEnd)
	if err != nil {
		log.Error(err.Error())
		return nil, customerror.Error(customerror.ErrCodeGeneric, customerror.ErrMessageGeneric)
	}

	response, err := buildOrderListResponse(orders, totalOrder, param.Page, param.Limit)
	if err != nil {
		log.Error(err.Error())
		return nil, customerror.Error(customerror.ErrCodeGeneric, customerror.ErrMessageGeneric)
	}

	return response, nil
}

func orderListToOrderListResponse(orders []repository.OrderList) ([]OrderListResponse, error) {
	result := []OrderListResponse{}
	for _, order := range orders {
		deliveryAmount := "-"

		if order.DeliveryAmount.Valid {
			deliveryAmount = fmt.Sprintf("$%.2f", order.DeliveryAmount.Float64)
		}

		location, err := time.LoadLocation(constant.MELBOURNE_LOCATION)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		result = append(result, OrderListResponse{
			OrderName:       order.OrderName,
			CustomerCompany: order.CustomerCompany,
			CustomerName:    order.CustomerName,
			OrderDate:       order.OrderDate.In(location),
			DeliveryAmount:  deliveryAmount,
			TotalAmount:     fmt.Sprintf("$%.2f", order.TotalAmount),
		})
	}

	return result, nil
}

func buildOrderListResponse(orders []repository.OrderList, totalOrder *repository.TotalOrder, currentPage int, limit int) (*OrderResponse, error) {
	data, err := orderListToOrderListResponse(orders)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	totalData := totalOrder.TotalData
	totalPage := int64(math.Ceil(float64(totalData) / float64(limit)))
	totalAmount := fmt.Sprintf("$%.2f", totalOrder.TotalAmount)
	isNext := int64(currentPage) != totalPage

	return &OrderResponse{
		CurrentPage: currentPage,
		Data:        data,
		TotalData:   totalData,
		TotalAmount: totalAmount,
		IsNext:      isNext,
		TotalPage:   totalPage,
	}, nil
}
