package repository

import (
	"card_manage/internal/model"
	"database/sql"
	"time"
)

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}

// CreateCard inserts a new card into the database.
func (r *CardRepository) CreateCard(card *model.Card) (int64, error) {
	query := `INSERT INTO cards (store_id, name, series, rarity, card_number, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	
	card.CreatedAt = time.Now()
	card.UpdatedAt = time.Now()

	var cardID int64
	err := r.db.QueryRow(
		query,
		card.StoreID,
		card.Name,
		card.Series,
		card.Rarity,
		card.CardNumber,
		card.CreatedAt,
		card.UpdatedAt,
	).Scan(&cardID)

	if err != nil {
		return 0, err
	}
	return cardID, nil
}

// GetCardByID retrieves a single card by its ID.
func (r *CardRepository) GetCardByID(cardID int64) (*model.Card, error) {
	query := `SELECT id, store_id, name, series, rarity, card_number, created_at, updated_at 
			  FROM cards WHERE id = $1`
	
	card := &model.Card{}
	err := r.db.QueryRow(query, cardID).Scan(
		&card.ID,
		&card.StoreID,
		&card.Name,
		&card.Series,
		&card.Rarity,
		&card.CardNumber,
		&card.CreatedAt,
		&card.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return card, nil
}

// UpdateCard updates an existing card in the database.
func (r *CardRepository) UpdateCard(card *model.Card) error {
	query := `UPDATE cards 
			  SET name = $1, series = $2, rarity = $3, card_number = $4, updated_at = $5
			  WHERE id = $6`
	
	card.UpdatedAt = time.Now()
	
	_, err := r.db.Exec(
		query,
		card.Name,
		card.Series,
		card.Rarity,
		card.CardNumber,
		card.UpdatedAt,
		card.ID,
	)
	return err
}

// DeleteCard removes a card from the database.
func (r *CardRepository) DeleteCard(cardID int64) error {
	query := `DELETE FROM cards WHERE id = $1`
	_, err := r.db.Exec(query, cardID)
	return err
}

// ListCardsByStore retrieves a list of cards for a specific store.
func (r *CardRepository) ListCardsByStore(storeID int64) ([]model.Card, error) {
	query := `SELECT id, store_id, name, series, rarity, card_number, created_at, updated_at
			  FROM cards WHERE store_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []model.Card
	for rows.Next() {
		var card model.Card
		if err := rows.Scan(
			&card.ID,
			&card.StoreID,
			&card.Name,
			&card.Series,
			&card.Rarity,
			&card.CardNumber,
			&card.CreatedAt,
			&card.UpdatedAt,
		); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}
