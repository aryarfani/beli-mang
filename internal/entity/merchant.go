package entity

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	ID        uuid.UUID `json:"userId"`
	Name      string    `json:"name"`
	Category  string    `json:"merchantCategory"`
	ImageUrl  string    `json:"imageUrl" db:"image_url"`
	Latitude  float64   `json:"lat"`
	Longitude float64   `json:"long"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
