package service

import (
	"card_manage/internal/model"
	"card_manage/internal/repository"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrConsignmentAlreadySold = errors.New("consignment has already been sold")
)

type TransactionService struct {
	repo            *repository.TransactionRepository
	consignmentRepo *repository.ConsignmentRepository
	storeRepo       *repository.StoreRepository
	db              *sql.DB
}

func NewTransactionService(
	repo *repository.TransactionRepository,
	consignmentRepo *repository.ConsignmentRepository,
	storeRepo *repository.StoreRepository,
	db *sql.DB,
) *TransactionService {
	return &TransactionService{
		repo:            repo,
		consignmentRepo: consignmentRepo,
		storeRepo:       storeRepo,
		db:              db,
	}
}

// CreateTransaction creates a new transaction and updates the consignment status within a DB transaction.
func (s *TransactionService) CreateTransaction(storeUserID, consignmentID int64, price float64, paymentMethod model.PaymentMethod) (*model.Transaction, error) {
	// 1. Get consignment and verify ownership and status
	consignment, err := s.consignmentRepo.GetConsignmentByID(consignmentID)
	if err != nil {
		return nil, fmt.Errorf("error getting consignment: %w", err)
	}
	if consignment == nil {
		return nil, ErrConsignmentNotFound
	}
	if consignment.Status == model.StatusSold || consignment.Status == model.StatusCleared {
		return nil, ErrConsignmentAlreadySold
	}

	store, err := s.storeRepo.GetStoreByUserID(storeUserID)
	if err != nil {
		return nil, fmt.Errorf("error getting store: %w", err)
	}
	if store == nil || store.ID != consignment.StoreID {
		return nil, ErrForbidden
	}

	// 2. Determine commission rate
	var commissionRate float64
	if paymentMethod == model.PaymentMethodCash {
		commissionRate = store.CommissionCash
	} else {
		commissionRate = store.CommissionCredit
	}

	// 3. Create transaction object
	newTxModel := &model.Transaction{
		ConsignmentID:  consignmentID,
		StoreID:        store.ID,
		Price:          price,
		PaymentMethod:  paymentMethod,
		CommissionRate: commissionRate,
	}

	// 4. Execute in a DB transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin db transaction: %w", err)
	}
	defer tx.Rollback() // Rollback is a no-op if tx has been committed.

	// Create the transaction record
	// Note: We need to adapt repository methods to accept a transaction (tx)
	// For now, we will create new methods for this purpose.
	txID, err := s.repo.CreateTransactionInTx(tx, newTxModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}
	newTxModel.ID = txID

	// Update the consignment status to SOLD
	if err := s.consignmentRepo.UpdateConsignmentStatusInTx(tx, consignmentID, model.StatusSold); err != nil {
		return nil, fmt.Errorf("failed to update consignment status: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit db transaction: %w", err)
	}

	return newTxModel, nil
}
