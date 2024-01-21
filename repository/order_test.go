package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"order-backend/constant"
	"time"
)

var _ = Describe("OrderRepository", func() {
	var repo *Order
	var mock sqlmock.Sqlmock

	BeforeEach(func() {
		var dbMock *sql.DB

		dbMock, mock, _ = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		db := sqlx.NewDb(dbMock, "oracle")

		repo = NewOrderRepository(db)
	})

	AfterEach(func() {
		err := mock.ExpectationsWereMet() // make sure all expectations were met
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("FindList", func() {
		var query = `select o.order_name                                               as order_name,
       cc.company_name                                            as customer_company,
       c.name                                                     as customer_name,
       o.created_at                                               as order_date,
       sum(COALESCE(oi.price_per_unit, 0) * d.delivered_quantity) as delivered_amount,
       sum(COALESCE(oi.price_per_unit, 0) * oi.quantity)          as total_amount
from orders o
         join customers c on c.user_id = o.customer_id
         join customer_companies cc on cc.company_id = c.company_id
         join order_items oi on o.id = oi.order_id
         left join deliveries d on oi.id = d.order_item_id
where 
   (
    ($1::TEXT is null or o.order_name ilike ('%' || $1 || '%'))
        or
    ($1::TEXT is null or oi.product ilike ('%' || $1 || '%'))
    )
  and (
        ($2::TIMESTAMPTZ IS NOT NULL AND $3::TIMESTAMPTZ IS NOT NULL AND o.created_at BETWEEN $2::TIMESTAMPTZ AND $3::TIMESTAMPTZ)
        OR
        ($2::TIMESTAMPTZ IS NOT NULL AND $3::TIMESTAMPTZ IS NULL AND o.created_at >= $2::TIMESTAMPTZ)
        OR
        ($2::TIMESTAMPTZ IS NULL AND $3::TIMESTAMPTZ IS NOT NULL AND o.created_at < $3::TIMESTAMPTZ)
        OR
        ($2::TIMESTAMPTZ IS NULL AND $3::TIMESTAMPTZ IS NULL)
    )
group by o.order_name, cc.company_name, c.name, o.created_at
order by order_date asc
limit $4 offset $5`

		orderMock := OrderList{
			OrderName:       "test order name",
			CustomerCompany: "test customer company",
			CustomerName:    "test customer name",
			OrderDate:       time.Now(),
			DeliveryAmount: sql.NullFloat64{
				Float64: 10.0,
				Valid:   true,
			},
			TotalAmount: 10.0,
		}

		orderColumn := []string{
			"order_name",
			"customer_company",
			"customer_name",
			"order_date",
			"delivered_amount",
			"total_amount",
		}

		loc, _ := time.LoadLocation(constant.MELBOURNE_LOCATION)
		startDate := time.Date(2020, 11, 20, 0, 0, 0, 0, loc)
		endDate := time.Date(2020, 12, 20, 0, 0, 0, 0, loc)

		orderRow := sqlmock.NewRows(orderColumn).AddRow(
			orderMock.OrderName,
			orderMock.CustomerCompany,
			orderMock.CustomerName,
			orderMock.OrderDate,
			orderMock.DeliveryAmount,
			orderMock.TotalAmount,
		)

		DescribeTable("when repository returns an error", func(repoError error) {
			mock.ExpectQuery(query).
				WillReturnError(repoError)
			result, err := repo.FindList(context.Background(), 10, 1, "test", &startDate, &endDate)
			Expect(err).To(Equal(repoError))
			Expect(result).To(BeNil())

		},
			Entry("repository returns an error", errors.New("repository error")),
		)
		DescribeTable("when repository returns success", func(expectedResults []OrderList, repoError error) {
			mock.ExpectQuery(query).
				WillReturnRows(orderRow).
				WillReturnError(repoError)

			result, err := repo.FindList(context.Background(), 10, 1, "test", &startDate, &endDate)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedResults))
		},
			Entry("repository returns data account", []OrderList{orderMock}, nil),
		)
	})

	Describe("GetTotalOrder", func() {
		var query = `select 
    sum(COALESCE(oi.price_per_unit, 0) * oi.quantity) as total_amount,
    count(distinct o.id) as total_data
from orders o
         join order_items oi on o.id = oi.order_id
where 
    (
    ($1::TEXT is null or o.order_name ilike ('%' || $1 || '%'))
        or
    ($1::TEXT is null or oi.product ilike ('%' || $1 || '%'))
    )
  and (
        ($2::TIMESTAMPTZ IS NOT NULL AND $3::TIMESTAMPTZ IS NOT NULL AND o.created_at BETWEEN $2::TIMESTAMPTZ AND $3::TIMESTAMPTZ)
        OR
        ($2::TIMESTAMPTZ IS NOT NULL AND $3::TIMESTAMPTZ IS NULL AND o.created_at >= $2::TIMESTAMPTZ)
        OR
        ($2::TIMESTAMPTZ IS NULL AND $3::TIMESTAMPTZ IS NOT NULL AND o.created_at < $3::TIMESTAMPTZ)
        OR
        ($2::TIMESTAMPTZ IS NULL AND $3::TIMESTAMPTZ IS NULL)
    )`

		orderMock := TotalOrder{
			TotalData:   10,
			TotalAmount: 10.0,
		}

		orderColumn := []string{
			"total_amount",
			"total_data",
		}

		loc, _ := time.LoadLocation(constant.MELBOURNE_LOCATION)
		startDate := time.Date(2020, 11, 20, 0, 0, 0, 0, loc)
		endDate := time.Date(2020, 12, 20, 0, 0, 0, 0, loc)

		orderRow := sqlmock.NewRows(orderColumn).AddRow(
			orderMock.TotalAmount,
			orderMock.TotalData,
		)

		DescribeTable("when repository returns an error", func(repoError error) {
			mock.ExpectQuery(query).
				WillReturnError(repoError)
			result, err := repo.FindTotalOrder(context.Background(), "test", &startDate, &endDate)
			Expect(err).To(Equal(repoError))
			Expect(result).To(BeNil())

		},
			Entry("repository returns an error", errors.New("repository error")),
		)
		DescribeTable("when repository returns success", func(expectedResults *TotalOrder, repoError error) {
			mock.ExpectQuery(query).
				WillReturnRows(orderRow).
				WillReturnError(repoError)

			result, err := repo.FindTotalOrder(context.Background(), "test", &startDate, &endDate)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedResults))
		},
			Entry("repository returns data account", &orderMock, nil),
		)
	})
})
