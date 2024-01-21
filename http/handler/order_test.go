package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	mocks "order-backend/mock/http/handler"
	"order-backend/order"
	"order-backend/pkg/utils"
	"time"
)

var _ = Describe("Handler", func() {
	Describe("Order List", func() {
		var orderServiceMock *mocks.OrderService
		var handler *OrderHandler

		BeforeEach(func() {
			orderServiceMock = &mocks.OrderService{}
			handler = NewOrderHandler(orderServiceMock)
		})

		location, _ := time.LoadLocation("Australia/Melbourne")

		date := time.Date(2020, 1, 1, 0, 0, 0, 0, location)
		orderList := order.OrderResponse{
			TotalData:   1,
			TotalAmount: "$100.00",
			TotalPage:   1,
			CurrentPage: 1,
			IsNext:      false,
			Data: []order.OrderListResponse{
				{
					OrderName:       "test order name",
					CustomerCompany: "test customer company",
					CustomerName:    "test customer name",
					OrderDate:       date,
					DeliveryAmount:  "$100.00",
					TotalAmount:     "$100.00",
				},
			},
		}

		startDate, _ := utils.FormatStartDate("2020-01-01")
		endDate, _ := utils.FormatEndDate("2020-01-01")

		DescribeTable("get order list success", func(query string, param order.Param, expectedStatusCode int, expectedResponse string) {
			orderServiceMock.On("GetListOrder", mock.Anything, param).Return(&orderList, nil)
			app := fiber.New()

			app.Get("/order", handler.GetList)

			req := httptest.NewRequest(http.MethodGet, "/order?"+query, nil)

			resp, err := app.Test(req, -1)
			Expect(err).NotTo(HaveOccurred())

			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(expectedStatusCode))
			Expect(string(body)).To(Equal(expectedResponse))
		},
			Entry("when send full param it should be return success",
				"date_end=2020-01-01&date_start=2020-01-01&keyword=test&limit=10&page=1",
				order.Param{
					Limit:     10,
					Page:      1,
					Keyword:   "test",
					DateStart: startDate,
					DateEnd:   endDate,
				}, http.StatusOK,
				`{"success":true,"data":{"total_data":1,"total_amount":"$100.00","total_page":1,"current_page":1,"next_page":false,"data":[{"order_name":"test order name","customer_company":"test customer company","customer_name":"test customer name","order_date":"2020-01-01T00:00:00+11:00","delivered_amount":"$100.00","total_amount":"$100.00"}]},"error":null}`,
			),
			Entry("when send without limit param it should be return success",
				"date_end=2020-01-01&date_start=2020-01-01&keyword=test&page=1",
				order.Param{
					Limit:     5,
					Page:      1,
					Keyword:   "test",
					DateStart: startDate,
					DateEnd:   endDate,
				}, http.StatusOK,
				`{"success":true,"data":{"total_data":1,"total_amount":"$100.00","total_page":1,"current_page":1,"next_page":false,"data":[{"order_name":"test order name","customer_company":"test customer company","customer_name":"test customer name","order_date":"2020-01-01T00:00:00+11:00","delivered_amount":"$100.00","total_amount":"$100.00"}]},"error":null}`,
			),
			Entry("when send without page param it should be return success",
				"date_end=2020-01-01&date_start=2020-01-01&keyword=test&limit=10",
				order.Param{
					Limit:     10,
					Page:      1,
					Keyword:   "test",
					DateStart: startDate,
					DateEnd:   endDate,
				}, http.StatusOK,
				`{"success":true,"data":{"total_data":1,"total_amount":"$100.00","total_page":1,"current_page":1,"next_page":false,"data":[{"order_name":"test order name","customer_company":"test customer company","customer_name":"test customer name","order_date":"2020-01-01T00:00:00+11:00","delivered_amount":"$100.00","total_amount":"$100.00"}]},"error":null}`,
			),
			Entry("when send without date start param it should be return success",
				"date_end=2020-01-01&keyword=test&limit=10&page=1",
				order.Param{
					Limit:   10,
					Page:    1,
					Keyword: "test",
				}, http.StatusOK,
				`{"success":true,"data":{"total_data":1,"total_amount":"$100.00","total_page":1,"current_page":1,"next_page":false,"data":[{"order_name":"test order name","customer_company":"test customer company","customer_name":"test customer name","order_date":"2020-01-01T00:00:00+11:00","delivered_amount":"$100.00","total_amount":"$100.00"}]},"error":null}`,
			),
		)

		DescribeTable("get order list error", func(handleMock func(), query string, expectedStatusCode int, expectedResponse string) {
			handleMock()

			app := fiber.New()

			app.Get("/order", handler.GetList)

			req := httptest.NewRequest(http.MethodGet, "/order?"+query, nil)

			resp, err := app.Test(req, -1)
			Expect(err).NotTo(HaveOccurred())

			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(expectedStatusCode))
			Expect(string(body)).To(Equal(expectedResponse))
		},
			Entry("when send invalid page param it should be return error",
				func() {},
				"date_end=2020-01-01&date_start=2020-01-01&keyword=test&limit=10&page=test",
				http.StatusBadRequest,
				`{"success":false,"data":null,"error":{"code":"error_code_bad_request","message":"invalid request"}}`,
			),
			Entry("when send invalid limit param it should be return error",
				func() {},
				"date_end=2020-01-01&date_start=2020-01-01&keyword=test&limit=test&page=1",
				http.StatusBadRequest,
				`{"success":false,"data":null,"error":{"code":"error_code_bad_request","message":"invalid request"}}`,
			),
			Entry("when send invalid format date start param it should be return error",
				func() {},
				"date_end=2020-01-01&date_start=2020-01-01T15:00:00&keyword=test&limit=10&page=1",
				http.StatusBadRequest,
				`{"success":false,"data":null,"error":{"code":"error_invalid_date_format","message":"invalid date format"}}`,
			),
			Entry("when send invalid format date end param it should be return error",
				func() {},
				"date_end=2020-01-01T15:00:00&date_start=2020-01-01&keyword=test&limit=10&page=1",
				http.StatusBadRequest,
				`{"success":false,"data":null,"error":{"code":"error_invalid_date_format","message":"invalid date format"}}`,
			),
			Entry("when service error it should be return error",
				func() {
					orderServiceMock.On("GetListOrder", mock.Anything, order.Param{
						Limit:     10,
						Page:      1,
						Keyword:   "test",
						DateStart: startDate,
						DateEnd:   endDate,
					}).Return(nil, errors.New("unexpected error"))
				},
				"date_end=2020-01-01&date_start=2020-01-01&keyword=test&limit=10&page=1",
				http.StatusInternalServerError,
				"{\"success\":false,\"data\":null,\"error\":{\"code\":\"error_code_generic\",\"message\":\"unexpected error\"}}",
			),
		)
	})
})
