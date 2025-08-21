package repository

import (
	"card_manage/internal/model"
	"database/sql"
	"time"
)

type ConsignmentRepository struct {
	db *sql.DB
}

func NewConsignmentRepository(db *sql.DB) *ConsignmentRepository {
	return &ConsignmentRepository{db: db}
}

// CreateConsignment inserts a new consignment into the database.
func (r *ConsignmentRepository) CreateConsignment(consignment *model.Consignment) (int64, error) {
	query := `INSERT INTO consignments (player_id, store_id, card_id, quantity, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	consignment.CreatedAt = time.Now()
	consignment.UpdatedAt = time.Now()

	var consignmentID int64
	err := r.db.QueryRow(
		query,
		consignment.PlayerID,
		consignment.StoreID,
		consignment.CardID,
		consignment.Quantity,
		consignment.Status,
		consignment.CreatedAt,
		consignment.UpdatedAt,
	).Scan(&consignmentID)

	if err != nil {
		return 0, err
	}
	return consignmentID, nil
}

// GetConsignmentByID retrieves a single consignment by its ID.
func (r *ConsignmentRepository) GetConsignmentByID(id int64) (*model.Consignment, error) {
	query := `SELECT id, player_id, store_id, card_id, quantity, status, created_at, updated_at 
			  FROM consignments WHERE id = $1`
	
	consignment := &model.Consignment{}
	err := r.db.QueryRow(query, id).Scan(
		&consignment.ID,
		&consignment.PlayerID,
		&consignment.StoreID,
		&consignment.CardID,
		&consignment.Quantity,
		&consignment.Status,
		&consignment.CreatedAt,
		&consignment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return consignment, nil
}

// ListConsignmentsByPlayer retrieves all consignments for a specific player.
func (r *ConsignmentRepository) ListConsignmentsByPlayer(playerID int64) ([]model.Consignment, error) {
	query := `SELECT id, player_id, store_id, card_id, quantity, status, created_at, updated_at 
			  FROM consignments WHERE player_id = $1 ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consignments []model.Consignment
	for rows.Next() {
		var c model.Consignment
		if err := rows.Scan(&c.ID, &c.PlayerID, &c.StoreID, &c.CardID, &c.Quantity, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		consignments = append(consignments, c)
	}
	return consignments, nil
}

// ListConsignmentsByStore retrieves all consignments for a specific store.
func (r *ConsignmentRepository) ListConsignmentsByStore(storeID int64) ([]model.Consignment, error) {
	query := `SELECT id, player_id, store_id, card_id, quantity, status, created_at, updated_at 
			  FROM consignments WHERE store_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consignments []model.Consignment
	for rows.Next() {
		var c model.Consignment
		if err := rows.Scan(&c.ID, &c.PlayerID, &c.StoreID, &c.CardID, &c.Quantity, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		consignments = append(consignments, c)
	}
	return consignments, nil
}

// UpdateConsignmentStatus updates the status of a consignment.
func (r *ConsignmentRepository) UpdateConsignmentStatus(id int64, status model.ConsignmentStatus) error {
	query := `UPDATE consignments SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

// UpdateConsignmentStatusInTx updates the status of a consignment within a specific transaction.
func (r *ConsignmentRepository) UpdateConsignmentStatusInTx(tx *sql.Tx, id int64, status model.ConsignmentStatus) error {
	query := `UPDATE consignments SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := tx.Exec(query, status, time.Now(), id)
	return err
}
