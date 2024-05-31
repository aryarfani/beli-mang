package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	UserId       uuid.UUID   `json:"userId" db:"user_id"`
	MerchantId   uuid.UUID   `json:"merchantId" db:"merchant_id"`
	EstimationId uuid.UUID   `json:"estimationId" db:"estimation_id"`
	CreatedAt    time.Time   `json:"createdAt" db:"created_at"`
	OrderItems   []OrderItem `json:"orderItems" db:"-"`
}
