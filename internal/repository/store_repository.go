package repository

import (
	"card_manage/internal/model"
	"database/sql"
	"time"
)

type StoreRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

// CreateStore inserts a new store into the database.
func (r *StoreRepository) CreateStore(store *model.Store) (int64, error) {
	query := `INSERT INTO stores (user_id, name, commission_cash, commission_credit, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	store.CreatedAt = time.Now()
	store.UpdatedAt = time.Now()

	var storeID int64
	err := r.db.QueryRow(
		query,
		store.UserID,
		store.Name,
		store.CommissionCash,
		store.CommissionCredit,
		store.CreatedAt,
		store.UpdatedAt,
	).Scan(&storeID)

	if err != nil {
		return 0, err
	}

	return storeID, nil
}

// GetStoreByUserID retrieves a store from the database by its owner's user ID.
func (r *StoreRepository) GetStoreByUserID(userID int64) (*model.Store, error) {
	query := `SELECT id, user_id, name, commission_cash, commission_credit, created_at, updated_at 
			  FROM stores WHERE user_id = $1`

	store := &model.Store{}
	err := r.db.QueryRow(query, userID).Scan(
		&store.ID,
		&store.UserID,
		&store.Name,
		&store.CommissionCash,
		&store.CommissionCredit,
		&store.CreatedAt,
		&store.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No store found is not an application error
		}
		return nil, err
	}

	return store, nil
}
