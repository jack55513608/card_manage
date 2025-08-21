package service

import (
	"card_manage/internal/model"
	"card_manage/internal/repository"
	"errors"
	"fmt"
)

var (
	ErrCardNotFound      = errors.New("card not found")
	ErrStoreNotFound     = errors.New("store not found for the current user")
	ErrForbidden         = errors.New("user is not allowed to perform this action")
)

type CardService struct {
	cardRepo  *repository.CardRepository
	storeRepo *repository.StoreRepository
}

func NewCardService(cardRepo *repository.CardRepository, storeRepo *repository.StoreRepository) *CardService {
	return &CardService{
		cardRepo:  cardRepo,
		storeRepo: storeRepo,
	}
}

// CreateCard creates a new card for the store associated with the given userID.
func (s *CardService) CreateCard(userID int64, name, series, rarity, cardNumber string) (*model.Card, error) {
	store, err := s.storeRepo.GetStoreByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error finding store: %w", err)
	}
	if store == nil {
		return nil, ErrStoreNotFound
	}

	newCard := &model.Card{
		StoreID:    store.ID,
		Name:       name,
		Series:     series,
		Rarity:     rarity,
		CardNumber: cardNumber,
	}

	cardID, err := s.cardRepo.CreateCard(newCard)
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}
	newCard.ID = cardID

	return newCard, nil
}

// GetCard checks if the user has permission and retrieves a card.
func (s *CardService) GetCard(userID, cardID int64) (*model.Card, error) {
	card, err := s.cardRepo.GetCardByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("error getting card: %w", err)
	}
	if card == nil {
		return nil, ErrCardNotFound
	}

	// Verify ownership
	if err := s.verifyStoreOwnership(userID, card.StoreID); err != nil {
		return nil, err
	}

	return card, nil
}

// ListCardsByCurrentUser lists all cards for the current user's store.
func (s *CardService) ListCardsByCurrentUser(userID int64) ([]model.Card, error) {
	store, err := s.storeRepo.GetStoreByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error finding store: %w", err)
	}
	if store == nil {
		return []model.Card{}, nil // Return empty slice if no store
	}

	return s.cardRepo.ListCardsByStore(store.ID)
}

// UpdateCard handles the logic for updating a card.
func (s *CardService) UpdateCard(userID, cardID int64, name, series, rarity, cardNumber string) (*model.Card, error) {
	// First, get the existing card
	card, err := s.cardRepo.GetCardByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("error getting card for update: %w", err)
	}
	if card == nil {
		return nil, ErrCardNotFound
	}

	// Verify ownership
	if err := s.verifyStoreOwnership(userID, card.StoreID); err != nil {
		return nil, err
	}

	// Update fields
	card.Name = name
	card.Series = series
	card.Rarity = rarity
	card.CardNumber = cardNumber

	if err := s.cardRepo.UpdateCard(card); err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}

	return card, nil
}

// DeleteCard handles the logic for deleting a card.
func (s *CardService) DeleteCard(userID, cardID int64) error {
	// First, get the existing card to verify ownership
	card, err := s.cardRepo.GetCardByID(cardID)
	if err != nil {
		return fmt.Errorf("error getting card for deletion: %w", err)
	}
	if card == nil {
		return ErrCardNotFound
	}

	// Verify ownership
	if err := s.verifyStoreOwnership(userID, card.StoreID); err != nil {
		return err
	}

	// Delete the card
	return s.cardRepo.DeleteCard(cardID)
}


// verifyStoreOwnership is a helper function to check if the user owns the store.
func (s *CardService) verifyStoreOwnership(userID, storeID int64) error {
	store, err := s.storeRepo.GetStoreByUserID(userID)
	if err != nil {
		return fmt.Errorf("error finding store for verification: %w", err)
	}
	if store == nil || store.ID != storeID {
		return ErrForbidden
	}
	return nil
}
