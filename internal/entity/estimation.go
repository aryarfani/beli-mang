package entity

import (
	"time"

	"github.com/google/uuid"
)

type Estimation struct {
	ID                             uuid.UUID `json:"id"`
	TotalPrice                     int       `json:"totalPrice" db:"total_price"`
	EstimatedDeliveryTimeInMinutes int       `json:"estimatedAt" db:"estimated_delivery_time_in_minutes"`
	IsOrdered                      bool      `json:"isOrdered" db:"is_ordered"`
	CreatedAt                      time.Time `json:"createdAt" db:"created_at"`
}
