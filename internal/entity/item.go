package entity

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID         uuid.UUID `json:"itemId"`
	MerchantId string    `json:"merchantId" db:"merchant_id"`
	Name       string    `json:"name"`
	Category   string    `json:"productCategory" db:"category"`
	Price      int       `json:"price"`
	ImageUrl   string    `json:"imageUrl" db:"image_url"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

var ItemCategories = []string{"Beverage", "Food", "Snack", "Condiments", "Additions"}
