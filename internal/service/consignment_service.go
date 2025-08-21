package service

import (
	"card_manage/internal/model"
	"card_manage/internal/repository"
	"fmt"
	"errors"
)

var (
	ErrConsignmentNotFound = errors.New("consignment not found")
	ErrInvalidCardForStore = errors.New("the selected card does not belong to the selected store")
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

// CreateConsignment allows a player to create a new consignment request.
func (s *ConsignmentService) CreateConsignment(playerID, storeID, cardID int64, quantity int) (*model.Consignment, error) {
	// 1. Validate that the card belongs to the store
	card, err := s.cardRepo.GetCardByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("error getting card: %w", err)
	}
	if card == nil || card.StoreID != storeID {
		return nil, ErrInvalidCardForStore
	}

	// 2. Create the consignment
	newConsignment := &model.Consignment{
		PlayerID: playerID,
		StoreID:  storeID,
		CardID:   cardID,
		Quantity: quantity,
		Status:   model.StatusPending,
	}

	consignmentID, err := s.consignmentRepo.CreateConsignment(newConsignment)
	if err != nil {
		return nil, fmt.Errorf("failed to create consignment: %w", err)
	}
	newConsignment.ID = consignmentID

	return newConsignment, nil
}

// ListConsignmentsForUser lists consignments based on user role.
// Players see their own consignments. Stores see consignments made to them.
func (s *ConsignmentService) ListConsignmentsForUser(userID int64, userRole string) ([]model.Consignment, error) {
	if userRole == "PLAYER" {
		return s.consignmentRepo.ListConsignmentsByPlayer(userID)
	}

	if userRole == "STORE" {
		store, err := s.storeRepo.GetStoreByUserID(userID)
		if err != nil {
			return nil, fmt.Errorf("error finding store for user: %w", err)
		}
		if store == nil {
			return []model.Consignment{}, nil // Store user might not have created a store profile yet
		}
		return s.consignmentRepo.ListConsignmentsByStore(store.ID)
	}

	return nil, ErrForbidden // Should not happen if RoleMiddleware is used correctly
}

// UpdateConsignmentStatus allows a store to update the status of a consignment.
func (s *ConsignmentService) UpdateConsignmentStatus(storeUserID, consignmentID int64, newStatus model.ConsignmentStatus) (*model.Consignment, error) {
	// 1. Get the consignment
	consignment, err := s.consignmentRepo.GetConsignmentByID(consignmentID)
	if err != nil {
		return nil, fmt.Errorf("error getting consignment: %w", err)
	}
	if consignment == nil {
		return nil, ErrConsignmentNotFound
	}

	// 2. Verify that the store owner is updating the consignment
	if err := s.verifyStoreOwnership(storeUserID, consignment.StoreID); err != nil {
		return nil, err
	}

	// 3. Update the status
	err = s.consignmentRepo.UpdateConsignmentStatus(consignmentID, newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update consignment status: %w", err)
	}
	consignment.Status = newStatus

	return consignment, nil
}

// Helper function from card_service - can be refactored into a shared service later
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
