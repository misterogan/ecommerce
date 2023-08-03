package model

type OrderItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func (oi *OrderItem) Insert() error {
	return nil
}
