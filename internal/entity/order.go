package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID   `json:"itemId"`
	UserId     uuid.UUID   `json:"userId"`
	MerchantId string      `json:"merchantId"`
	CreatedAt  time.Time   `json:"createdAt" db:"created_at"`
	OrderItems []OrderItem `json:"orderItems"`
}
