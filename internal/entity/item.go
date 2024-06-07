package entity

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID         uuid.UUID `json:"itemId"`
	MerchantId string    `json:"merchant_id" db:"merchant_id"`
	Name       string    `json:"name"`
	Category   string    `json:"category" db:"category"`
	Price      int       `json:"price"`
	ImageUrl   string    `json:"image_url" db:"image_url"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

var ItemCategories = []string{"Beverage", "Food", "Snack", "Condiments", "Additions"}
