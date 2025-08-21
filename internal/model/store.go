package model

import "time"

// Store corresponds to the "stores" table in the database.
type Store struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	Name             string    `json:"name"`
	CommissionCash   float64   `json:"commission_cash"`
	CommissionCredit float64   `json:"commission_credit"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
