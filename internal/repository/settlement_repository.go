package repository

import (
	"card_manage/internal/model"
	"database/sql"
	"time"
)

type SettlementRepository struct {
	db *sql.DB
}

func NewSettlementRepository(db *sql.DB) *SettlementRepository {
	return &SettlementRepository{db: db}
}

// CreateSettlement creates a new settlement request.
func (r *SettlementRepository) CreateSettlement(settlement *model.Settlement) (int64, error) {
	query := `INSERT INTO settlements (player_id, store_id, amount, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	
	settlement.CreatedAt = time.Now()
	settlement.UpdatedAt = time.Now()

	var settlementID int64
	err := r.db.QueryRow(
		query,
		settlement.PlayerID,
		settlement.StoreID,
		settlement.Amount,
		settlement.Status,
		settlement.CreatedAt,
		settlement.UpdatedAt,
	).Scan(&settlementID)

	if err != nil {
		return 0, err
	}
	return settlementID, nil
}

// ListSettlementsByPlayer retrieves all settlements for a specific player.
func (r *SettlementRepository) ListSettlementsByPlayer(playerID int64) ([]model.Settlement, error) {
	// Implementation similar to other List methods
	return nil, nil // Placeholder
}

// ListSettlementsByStore retrieves all settlements for a specific store.
func (r *SettlementRepository) ListSettlementsByStore(storeID int64) ([]model.Settlement, error) {
	// Implementation similar to other List methods
	return nil, nil // Placeholder
}

// UpdateSettlementStatus updates the status of a settlement.
func (r *SettlementRepository) UpdateSettlementStatus(id int64, status model.SettlementStatus) error {
	query := `UPDATE settlements SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

// GetUnsettledTransactions calculates the total amount from sold but uncleared consignments for a player at a specific store.
func (r *SettlementRepository) GetUnsettledTransactions(playerID, storeID int64) ([]model.Transaction, error) {
	query := `
		SELECT t.id, t.consignment_id, t.store_id, t.price, t.payment_method, t.commission_rate, t.created_at
		FROM transactions t
		JOIN consignments c ON t.consignment_id = c.id
		WHERE c.player_id = $1
		  AND c.store_id = $2
		  AND c.status = 'SOLD'`

	rows, err := r.db.Query(query, playerID, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var tx model.Transaction
		if err := rows.Scan(&tx.ID, &tx.ConsignmentID, &tx.StoreID, &tx.Price, &tx.PaymentMethod, &tx.CommissionRate, &tx.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}
