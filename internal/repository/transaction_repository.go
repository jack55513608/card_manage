package repository

import (
	"card_manage/internal/model"
	"database/sql"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction inserts a new transaction into the database.
func (r *TransactionRepository) CreateTransaction(tx *model.Transaction) (int64, error) {
	query := `INSERT INTO transactions (consignment_item_id, store_id, price, payment_method, commission_rate, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	tx.CreatedAt = time.Now()

	var transactionID int64
	err := r.db.QueryRow(
		query,
		tx.ConsignmentItemID,
		tx.StoreID,
		tx.Price,
		tx.PaymentMethod,
		tx.CommissionRate,
		tx.CreatedAt,
	).Scan(&transactionID)

	if err != nil {
		return 0, err
	}
	return transactionID, nil
}

// CreateTransactionInTx inserts a new transaction into the database within a specific transaction.
func (r *TransactionRepository) CreateTransactionInTx(tx *sql.Tx, transaction *model.Transaction) (int64, error) {
	query := `INSERT INTO transactions (consignment_item_id, store_id, price, payment_method, commission_rate, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	transaction.CreatedAt = time.Now()

	var transactionID int64
	err := tx.QueryRow(
		query,
		transaction.ConsignmentItemID,
		transaction.StoreID,
		transaction.Price,
		transaction.PaymentMethod,
		transaction.CommissionRate,
		transaction.CreatedAt,
	).Scan(&transactionID)

	if err != nil {
		return 0, err
	}
	return transactionID, nil
}
