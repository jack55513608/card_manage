package service

import (
	"card_manage/internal/model"
	"card_manage/internal/repository"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNoUnsettledTransactions = errors.New("no unsettled transactions found to create a settlement")
)

type SettlementService struct {
	repo            *repository.SettlementRepository
	consignmentRepo *repository.ConsignmentRepository
	storeRepo       *repository.StoreRepository
	db              *sql.DB
}

func NewSettlementService(
	repo *repository.SettlementRepository,
	consignmentRepo *repository.ConsignmentRepository,
	storeRepo *repository.StoreRepository,
	db *sql.DB,
) *SettlementService {
	return &SettlementService{
		repo:            repo,
		consignmentRepo: consignmentRepo,
		storeRepo:       storeRepo,
		db:              db,
	}
}

// CreateSettlement allows a player to request a settlement for a specific store.
func (s *SettlementService) CreateSettlement(playerID, storeID int64) (*model.Settlement, error) {
	// 1. Get all unsettled (SOLD) transactions for the player-store pair
	transactions, err := s.repo.GetUnsettledTransactions(playerID, storeID)
	if err != nil {
		return nil, fmt.Errorf("error getting unsettled transactions: %w", err)
	}
	if len(transactions) == 0 {
		return nil, ErrNoUnsettledTransactions
	}

	// 2. Calculate the total settlement amount
	var totalAmount float64
	var consignmentIDsToClear []int64
	for _, tx := range transactions {
		playerShare := tx.Price * (1 - (tx.CommissionRate / 100.0))
		totalAmount += playerShare
		consignmentIDsToClear = append(consignmentIDsToClear, tx.ConsignmentID)
	}

	// 3. Create the settlement object
	newSettlement := &model.Settlement{
		PlayerID: playerID,
		StoreID:  storeID,
		Amount:   totalAmount,
		Status:   model.StatusRequested,
	}

	// 4. Execute in a DB transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin db transaction: %w", err)
	}
	defer tx.Rollback()

	// Create the settlement record
	settlementID, err := s.repo.CreateSettlement(newSettlement)
	if err != nil {
		return nil, fmt.Errorf("failed to create settlement record: %w", err)
	}
	newSettlement.ID = settlementID

	// Update all related consignments to CLEARED
	for _, consignmentID := range consignmentIDsToClear {
		if err := s.consignmentRepo.UpdateConsignmentStatusInTx(tx, consignmentID, model.StatusCleared); err != nil {
			return nil, fmt.Errorf("failed to update consignment status to cleared: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit db transaction: %w", err)
	}

	return newSettlement, nil
}

// CompleteSettlement allows a store to mark a settlement as completed.
func (s *SettlementService) CompleteSettlement(storeUserID, settlementID int64) (*model.Settlement, error) {
	// Implementation for store to complete a settlement
	// 1. Get settlement
	// 2. Verify store ownership
	// 3. Update status
	return nil, errors.New("not implemented yet")
}
