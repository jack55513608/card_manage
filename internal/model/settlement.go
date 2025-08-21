package model

import "time"

// SettlementStatus represents the status of a settlement request.
type SettlementStatus string

const (
	StatusRequested SettlementStatus = "REQUESTED"
	StatusCompleted SettlementStatus = "COMPLETED"
)

// Settlement corresponds to the "settlements" table in the database.
type Settlement struct {
	ID        int64            `json:"id"`
	PlayerID  int64            `json:"player_id"`
	StoreID   int64            `json:"store_id"`
	Amount    float64          `json:"amount"`
	Status    SettlementStatus `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
