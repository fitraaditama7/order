package order

import (
	"time"
)

type OrderResponse struct {
	TotalData   int64               `json:"total_data"`
	TotalAmount string              `json:"total_amount"`
	TotalPage   int64               `json:"total_page"`
	CurrentPage int                 `json:"current_page"`
	IsNext      bool                `json:"next_page"`
	Data        []OrderListResponse `json:"data"`
}

type OrderListResponse struct {
	OrderName       string    `json:"order_name"`
	CustomerCompany string    `json:"customer_company"`
	CustomerName    string    `json:"customer_name"`
	OrderDate       time.Time `json:"order_date"`
	DeliveryAmount  string    `json:"delivered_amount"`
	TotalAmount     string    `json:"total_amount"`
}
