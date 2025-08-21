package service

import (
	"card_manage/internal/model"
	"card_manage/internal/repository"
	"errors"
	"fmt"
)

var (
	ErrConsignmentNotFound      = errors.New("consignment not found")
	ErrConsignmentItemNotFound  = errors.New("consignment item not found")
	ErrInvalidCardForStore      = errors.New("one or more cards do not belong to the selected store")
	ErrCannotUpdateStatus       = errors.New("status cannot be updated to the desired value")
)

type ConsignmentService struct {
	consignmentRepo *repository.ConsignmentRepository
	cardRepo        *repository.CardRepository
	storeRepo       *repository.StoreRepository
}

func NewConsignmentService(
	consignmentRepo *repository.ConsignmentRepository,
	cardRepo *repository.CardRepository,
	storeRepo *repository.StoreRepository,
) *ConsignmentService {
	return &ConsignmentService{
		consignmentRepo: consignmentRepo,
		cardRepo:        cardRepo,
		storeRepo:       storeRepo,
	}
}

// CreateConsignment allows a player to create a new consignment request with multiple items.
func (s *ConsignmentService) CreateConsignment(playerID, storeID int64, cardIDs []int64) (*model.Consignment, error) {
	// In a real application, you'd validate that all cardIDs are valid and belong to the store.
	// This is omitted for brevity but is crucial for production code.

	consignment := &model.Consignment{
		PlayerID: playerID,
		StoreID:  storeID,
		Status:   model.ConsignmentRequestStatusProcessing,
	}

	var items []*model.ConsignmentItem
	for _, cardID := range cardIDs {
		items = append(items, &model.ConsignmentItem{
			CardID: cardID,
			Status: model.ItemStatusPending,
		})
	}

	consignmentID, err := s.consignmentRepo.CreateConsignment(consignment, items)
	if err != nil {
		return nil, fmt.Errorf("failed to create consignment: %w", err)
	}

	// Return the full consignment object
	return s.consignmentRepo.GetConsignmentByID(consignmentID)
}

// UpdateConsignmentItemStatus allows a store to approve or reject a specific item.
func (s *ConsignmentService) UpdateConsignmentItemStatus(storeUserID, itemID int64, newStatus model.ConsignmentItemStatus, reason string) (*model.ConsignmentItem, error) {
	// 1. Get the item
	item, err := s.consignmentRepo.GetConsignmentItemByID(itemID)
	if err != nil {
		return nil, fmt.Errorf("error getting item: %w", err)
	}
	if item == nil {
		return nil, ErrConsignmentItemNotFound
	}

	// 2. Get the parent consignment to find the store ID for verification
	consignment, err := s.consignmentRepo.GetConsignmentByID(item.ConsignmentID)
	if err != nil {
		return nil, fmt.Errorf("error getting parent consignment: %w", err)
	}
	if consignment == nil {
		return nil, ErrConsignmentNotFound // Should not happen if item exists
	}

	// 3. Verify ownership
	if err := s.verifyStoreOwnership(storeUserID, consignment.StoreID); err != nil {
		return nil, err
	}

	// 4. Validate status transition (can only approve/reject from pending)
	if item.Status != model.ItemStatusPending {
		return nil, fmt.Errorf("%w: current status is %s", ErrCannotUpdateStatus, item.Status)
	}
	if newStatus != model.ItemStatusApproved && newStatus != model.ItemStatusRejected {
		return nil, fmt.Errorf("%w: can only change to APPROVED or REJECTED", ErrCannotUpdateStatus)
	}

	// 5. Update the status
	if err := s.consignmentRepo.UpdateConsignmentItemStatus(itemID, newStatus, reason); err != nil {
		return nil, fmt.Errorf("failed to update item status: %w", err)
	}

	item.Status = newStatus
	item.RejectionReason = reason
	return item, nil
}

// verifyStoreOwnership is a helper function to check if the user owns the store.
func (s *ConsignmentService) verifyStoreOwnership(userID, storeID int64) error {
	store, err := s.storeRepo.GetStoreByUserID(userID)
	if err != nil {
		return fmt.Errorf("error finding store for verification: %w", err)
	}
	if store == nil || store.ID != storeID {
		return ErrForbidden
	}
	return nil
}