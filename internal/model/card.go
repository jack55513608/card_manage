package model

import "time"

// Card corresponds to the "cards" table in the database.
type Card struct {
	ID         int64     `json:"id"`
	StoreID    int64     `json:"store_id"`
	Name       string    `json:"name"`
	Series     string    `json:"series,omitempty"`
	Rarity     string    `json:"rarity,omitempty"`
	CardNumber string    `json:"card_number,omitempty"`
	ImageURL   string    `json:"image_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
