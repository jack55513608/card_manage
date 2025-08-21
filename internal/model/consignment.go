package model

import "time"

// ConsignmentRequestStatus represents the status of a consignment request.
type ConsignmentRequestStatus string

const (
	ConsignmentRequestStatusProcessing ConsignmentRequestStatus = "PROCESSING"
	ConsignmentRequestStatusCompleted  ConsignmentRequestStatus = "COMPLETED"
)

// ConsignmentItemStatus represents the status of an individual item in a consignment.
type ConsignmentItemStatus string

const (
	ItemStatusPending   ConsignmentItemStatus = "PENDING"
	ItemStatusApproved  ConsignmentItemStatus = "APPROVED"
	ItemStatusRejected  ConsignmentItemStatus = "REJECTED"
	ItemStatusSold      ConsignmentItemStatus = "SOLD"
	ItemStatusCleared   ConsignmentItemStatus = "CLEARED"
)

// Consignment corresponds to the "consignments" table (a request).
type Consignment struct {
	ID        int64                    `json:"id"`
	PlayerID  int64                    `json:"player_id"`
	StoreID   int64                    `json:"store_id"`
	Status    ConsignmentRequestStatus `json:"status"`
	Items     []ConsignmentItem        `json:"items,omitempty"` // Used for API responses
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

// ConsignmentItem corresponds to the "consignment_items" table.
type ConsignmentItem struct {
	ID               int64                 `json:"id"`
	ConsignmentID    int64                 `json:"consignment_id"`
	CardID           int64                 `json:"card_id"`
	Status           ConsignmentItemStatus `json:"status"`
	RejectionReason  string                `json:"rejection_reason,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}