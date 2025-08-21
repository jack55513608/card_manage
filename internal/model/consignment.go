package model

import "time"

// ConsignmentStatus represents the status of a consignment.
type ConsignmentStatus string

const (
	StatusPending ConsignmentStatus = "PENDING"
	StatusListed  ConsignmentStatus = "LISTED"
	StatusSold    ConsignmentStatus = "SOLD"
	StatusCleared ConsignmentStatus = "CLEARED"
)

// Consignment corresponds to the "consignments" table in the database.
type Consignment struct {
	ID        int64             `json:"id"`
	PlayerID  int64             `json:"player_id"`
	StoreID   int64             `json:"store_id"`
	CardID    int64             `json:"card_id"`
	Quantity  int               `json:"quantity"`
	Status    ConsignmentStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
