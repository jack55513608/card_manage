package repository

import (
	"card_manage/internal/model"
	"database/sql"
	"fmt"
	"time"
)

type ConsignmentRepository struct {
	db *sql.DB
}

func NewConsignmentRepository(db *sql.DB) *ConsignmentRepository {
	return &ConsignmentRepository{db: db}
}

// CreateConsignment creates a new consignment request and its items in a single transaction.
func (r *ConsignmentRepository) CreateConsignment(consignment *model.Consignment, items []*model.ConsignmentItem) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on any error

	// 1. Create the parent consignment request
	consignmentQuery := `INSERT INTO consignments (player_id, store_id, status, created_at, updated_at)
					   VALUES ($1, $2, $3, $4, $5) RETURNING id`
	consignment.CreatedAt = time.Now()
	consignment.UpdatedAt = time.Now()

	var consignmentID int64
	err = tx.QueryRow(
		consignmentQuery,
		consignment.PlayerID,
		consignment.StoreID,
		consignment.Status,
		consignment.CreatedAt,
		consignment.UpdatedAt,
	).Scan(&consignmentID)
	if err != nil {
		return 0, fmt.Errorf("failed to create consignment request: %w", err)
	}

	// 2. Create the consignment items
	itemStmt, err := tx.Prepare(`INSERT INTO consignment_items (consignment_id, card_id, status, created_at, updated_at)
								VALUES ($1, $2, $3, $4, $5)`) 
	if err != nil {
		return 0, fmt.Errorf("failed to prepare item statement: %w", err)
	}
	defer itemStmt.Close()

	for _, item := range items {
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()
		_, err := itemStmt.Exec(consignmentID, item.CardID, item.Status, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return 0, fmt.Errorf("failed to create consignment item for card %d: %w", item.CardID, err)
		}
	}

	// 3. Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return consignmentID, nil
}

// GetConsignmentByID retrieves a consignment and all its items.
func (r *ConsignmentRepository) GetConsignmentByID(id int64) (*model.Consignment, error) {
	// Get the parent consignment
	consignment := &model.Consignment{}
	query := `SELECT id, player_id, store_id, status, created_at, updated_at FROM consignments WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&consignment.ID, &consignment.PlayerID, &consignment.StoreID, 
		&consignment.Status, &consignment.CreatedAt, &consignment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("error getting consignment: %w", err)
	}

	// Get the associated items
	rows, err := r.db.Query(`SELECT id, consignment_id, card_id, status, rejection_reason, created_at, updated_at 
							   FROM consignment_items WHERE consignment_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("error getting consignment items: %w", err)
	}
	defer rows.Close()

	var items []model.ConsignmentItem
	for rows.Next() {
		var item model.ConsignmentItem
		if err := rows.Scan(&item.ID, &item.ConsignmentID, &item.CardID, &item.Status, &item.RejectionReason, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning consignment item: %w", err)
		}
		items = append(items, item)
	}
	consignment.Items = items

	return consignment, nil
}

// GetConsignmentItemByID retrieves a single consignment item.
func (r *ConsignmentRepository) GetConsignmentItemByID(id int64) (*model.ConsignmentItem, error) {
	item := &model.ConsignmentItem{}
	query := `SELECT id, consignment_id, card_id, status, rejection_reason, created_at, updated_at 
			  FROM consignment_items WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&item.ID, &item.ConsignmentID, &item.CardID, &item.Status, 
		&item.RejectionReason, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("error getting consignment item: %w", err)
	}
	return item, nil
}

// UpdateConsignmentItemStatus updates the status and rejection reason of a specific item.
func (r *ConsignmentRepository) UpdateConsignmentItemStatus(id int64, status model.ConsignmentItemStatus, reason string) error {
	query := `UPDATE consignment_items SET status = $1, rejection_reason = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, status, reason, time.Now(), id)
	return err
}