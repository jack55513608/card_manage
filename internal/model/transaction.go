package model

import "time"

// PaymentMethod represents the payment method used for a transaction.
type PaymentMethod string

const (
	PaymentMethodCash   PaymentMethod = "CASH"
	PaymentMethodCredit PaymentMethod = "CREDIT"
)

// Transaction corresponds to the "transactions" table in the database.
type Transaction struct {
	ID             int64         `json:"id"`
	ConsignmentID  int64         `json:"consignment_id"`
	StoreID        int64         `json:"store_id"`
	Price          float64       `json:"price"`
	PaymentMethod  PaymentMethod `json:"payment_method"`
	CommissionRate float64       `json:"commission_rate"`
	CreatedAt      time.Time     `json:"created_at"`
}
