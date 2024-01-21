package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type Order struct {
	db *sqlx.DB
}

type OrderList struct {
	OrderName       string          `db:"order_name"`
	CustomerCompany string          `db:"customer_company"`
	CustomerName    string          `db:"customer_name"`
	OrderDate       time.Time       `db:"order_date"`
	DeliveryAmount  sql.NullFloat64 `db:"delivered_amount"`
	TotalAmount     float64         `db:"total_amount"`
}

type TotalOrder struct {
	TotalAmount float64 `db:"total_amount"`
	TotalData   int64   `db:"total_data"`
}

func NewOrderRepository(db *sqlx.DB) *Order {
	return &Order{db: db}
}

func (o *Order) FindList(ctx context.Context, limit, page int, keyword string, dateStart, dateEnd *time.Time) ([]OrderList, error) {
	var result = []OrderList{}

	offset := (page - 1) * limit

	args := []interface{}{
		keyword,
		dateStart,
		dateEnd,
		limit,
		offset,
	}

	query := `select o.order_name                                               as order_name,
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

	err := o.db.SelectContext(ctx, &result, query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (o *Order) FindTotalOrder(ctx context.Context, keyword string, dateStart, dateEnd *time.Time) (*TotalOrder, error) {
	var result TotalOrder

	args := []interface{}{
		keyword,
		dateStart,
		dateEnd,
	}

	query := `select 
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

	err := o.db.GetContext(ctx, &result, query, args...)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
