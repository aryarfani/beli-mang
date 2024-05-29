package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID        uuid.UUID `json:"id"`
	OrderId   uuid.UUID `json:"orderId"`
	ItemId    uuid.UUID `json:"itemId"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
