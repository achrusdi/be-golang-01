package models

type Payment struct {
	FromCustomerID string  `json:"from_customer_id"`
	ToCustomerID   string  `json:"to_customer_id"`
	Amount         float64 `json:"amount"`
	Timestamp      string  `json:"timestamp"`
}
