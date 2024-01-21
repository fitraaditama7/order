package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math"
	"order-backend/constant"
	mocks "order-backend/mock/order"
	"order-backend/pkg/customerror"
	"order-backend/repository"
	"time"
)

var _ = Describe("Order Service", func() {
	Describe("Order List", func() {
		var orderRepositoryMock *mocks.OrderRepository
		var services *service

		BeforeEach(func() {
			orderRepositoryMock = &mocks.OrderRepository{}
			services = NewOrderService(orderRepositoryMock)
		})

		ctx := context.Background()

		orderList := repository.OrderList{
			OrderName:       "test order name",
			CustomerCompany: "test customer company",
			CustomerName:    "test consumer name",
			OrderDate:       time.Now(),
			DeliveryAmount: sql.NullFloat64{
				Float64: 0,
				Valid:   true,
			},
			TotalAmount: 10000.0,
		}

		totalOrder := repository.TotalOrder{
			TotalAmount: 10000.0,
			TotalData:   1,
		}

		loc, _ := time.LoadLocation(constant.MELBOURNE_LOCATION)
		startDate := time.Date(2020, 11, 20, 0, 0, 0, 0, loc)
		endDate := time.Date(2020, 12, 20, 0, 0, 0, 0, loc)

		param := Param{
			Limit:     10,
			Page:      1,
			Keyword:   "TEST",
			DateStart: &startDate,
			DateEnd:   &endDate,
		}

		data, _ := orderListToOrderListResponse([]repository.OrderList{orderList})

		DescribeTable("when success get order list", func(
			orderList []repository.OrderList,
			totalOrder *repository.TotalOrder,
			expectedResult *OrderResponse,
		) {
			orderRepositoryMock.On("FindList", ctx,
				param.Limit,
				param.Page,
				param.Keyword,
				param.DateStart,
				param.DateEnd,
			).Return(orderList, nil)
			orderRepositoryMock.On("FindTotalOrder", ctx, param.Keyword, param.DateStart, param.DateEnd).
				Return(totalOrder, nil)

			result, err := services.GetListOrder(ctx, param)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedResult))
		},
			Entry("it should return success with data",
				[]repository.OrderList{orderList},
				&totalOrder,
				&OrderResponse{
					TotalData:   totalOrder.TotalData,
					TotalAmount: fmt.Sprintf("$%.2f", totalOrder.TotalAmount),
					TotalPage:   int64(math.Ceil(float64(totalOrder.TotalData) / float64(param.Limit))),
					CurrentPage: 1,
					IsNext:      false,
					Data:        data,
				},
			),
			Entry("it should return success with empty data",
				[]repository.OrderList{},
				&totalOrder,
				&OrderResponse{
					TotalData:   0,
					TotalPage:   1,
					CurrentPage: 1,
					IsNext:      false,
					Data:        []OrderListResponse{},
				},
			),
		)

		DescribeTable("error get order list", func(
			handleMock func(),
			expectedError error,
		) {
			handleMock()

			result, err := services.GetListOrder(ctx, param)
			Expect(err).To(Equal(expectedError))
			Expect(result).To(BeNil())
		},
			Entry("when failed to got unexpected error from repository it should return error",
				func() {
					orderRepositoryMock.On("FindList", ctx,
						param.Limit,
						param.Page,
						param.Keyword,
						param.DateStart,
						param.DateEnd,
					).Return(nil, errors.New("unexpected error"))
				},
				customerror.Error(customerror.ErrCodeGeneric, customerror.ErrMessageGeneric),
			),
			Entry("when failed to got unexpected error from repository it should return error",
				func() {
					orderRepositoryMock.On("FindList", ctx,
						param.Limit,
						param.Page,
						param.Keyword,
						param.DateStart,
						param.DateEnd,
					).Return([]repository.OrderList{orderList}, nil)
					orderRepositoryMock.On("FindTotalOrder", ctx, param.Keyword, param.DateStart, param.DateEnd).
						Return(nil, errors.New("unexpected error"))
				},
				customerror.Error(customerror.ErrCodeGeneric, customerror.ErrMessageGeneric),
			),
		)
	})
})
