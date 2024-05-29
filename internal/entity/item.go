package entity

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID        uuid.UUID `json:"itemId"`
	Name      string    `json:"name"`
	Category  string    `json:"productCategory"`
	Price     int       `json:"price"`
	ImageUrl  string    `json:"imageUrl" db:"image_url"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

var ItemCategories = []string{"Beverage", "Food", "Snack", "Condiments", "Additions"}
