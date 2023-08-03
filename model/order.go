package model

type Order struct {
	ID         int64       `json:"id"`
	CustomerID int64       `json:"customer_id"`
	OrderDate  string      `json:"order_date"`
	Status     string      `json:"status"`
	Items      []OrderItem `json:"items"`
}
