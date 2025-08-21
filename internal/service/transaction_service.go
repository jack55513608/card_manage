package service

import (
	"card_manage/internal/model"
	"card_manage/internal/repository"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrItemNotApproved = errors.New("consignment item is not approved for sale")
	ErrItemAlreadySold = errors.New("consignment item has already been sold")
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

// CreateTransaction creates a new transaction for a specific consignment item.
func (s *TransactionService) CreateTransaction(storeUserID, itemID int64, price float64, paymentMethod model.PaymentMethod) (*model.Transaction, error) {
	// 1. Get the consignment item and its parent consignment
	item, err := s.consignmentRepo.GetConsignmentItemByID(itemID)
	if err != nil {
		return nil, fmt.Errorf("error getting consignment item: %w", err)
	}
	if item == nil {
		return nil, ErrConsignmentItemNotFound
	}

	consignment, err := s.consignmentRepo.GetConsignmentByID(item.ConsignmentID)
	if err != nil {
		return nil, fmt.Errorf("error getting parent consignment: %w", err)
	}
	if consignment == nil {
		return nil, ErrConsignmentNotFound
	}

	// 2. Verify store ownership
	store, err := s.storeRepo.GetStoreByUserID(storeUserID)
	if err != nil {
		return nil, fmt.Errorf("error getting store: %w", err)
	}
	if store == nil || store.ID != consignment.StoreID {
		return nil, ErrForbidden
	}

	// 3. Check if the item can be sold
	if item.Status != model.ItemStatusApproved {
		return nil, ErrItemNotApproved
	}

	// 4. Determine commission rate
	var commissionRate float64
	if paymentMethod == model.PaymentMethodCash {
		commissionRate = store.CommissionCash
	} else {
		commissionRate = store.CommissionCredit
	}

	// 5. Create transaction object
	newTxModel := &model.Transaction{
		ConsignmentItemID: itemID,
		StoreID:           store.ID,
		Price:             price,
		PaymentMethod:     paymentMethod,
		CommissionRate:    commissionRate,
	}

	// 6. Execute in a DB transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin db transaction: %w", err)
	}
	defer tx.Rollback()

	// Create the transaction record
	txID, err := s.repo.CreateTransactionInTx(tx, newTxModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}
	newTxModel.ID = txID

	// Update the item status to SOLD
	if err := s.consignmentRepo.UpdateConsignmentItemStatus(itemID, model.ItemStatusSold, ""); err != nil {
		return nil, fmt.Errorf("failed to update item status: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit db transaction: %w", err)
	}

	return newTxModel, nil
}