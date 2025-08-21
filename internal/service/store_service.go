package service

import (
	"card_manage/internal/model"
	"card_manage/internal/repository"
	"fmt"
)

type StoreService struct {
	storeRepo *repository.StoreRepository
}

func NewStoreService(storeRepo *repository.StoreRepository) *StoreService {
	return &StoreService{storeRepo: storeRepo}
}

// CreateStore handles the business logic for creating a new store.
// It links the store to the user ID provided.
func (s *StoreService) CreateStore(userID int64, name string, commissionCash, commissionCredit float64) (*model.Store, error) {
	// In the future, we might add validation here, e.g., check if a user already has a store.

	newStore := &model.Store{
		UserID:           userID,
		Name:             name,
		CommissionCash:   commissionCash,
		CommissionCredit: commissionCredit,
	}

	storeID, err := s.storeRepo.CreateStore(newStore)
	if err != nil {
		// We might want to check for specific DB errors here, like a unique constraint violation on user_id
		return nil, fmt.Errorf("failed to create store in service: %w", err)
	}
	newStore.ID = storeID

	return newStore, nil
}
