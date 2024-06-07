package item

import (
	"beli-mang/internal/entity"
	"time"

	"github.com/google/uuid"
)

type CreateItemRequest struct {
	Name            string `json:"name" validate:"required,min=2,max=30"`
	ProductCategory string `json:"productCategory" validate:"required,oneof=Beverage Food Snack Condiments Additions"`
	Price           int    `json:"price" validate:"required,min=1"`
	ImageUrl        string `json:"imageUrl" validate:"required,xImageUrl"`
	MerchantId      string
}

type QueryItemsRequest struct {
	ItemID     string
	Limit      int
	Offset     int
	Name       string
	Category   string
	CreatedAt  string
	MerchantId string
}

type PaginatedQueryItemsResponse struct {
	Data []Item                `json:"data"`
	Meta entity.PaginationMeta `json:"meta"`
}

type Item struct {
	ID         uuid.UUID `json:"itemId"`
	MerchantId string    `json:"merchantId" db:"merchant_id"`
	Name       string    `json:"name"`
	Category   string    `json:"productCategory" db:"category"`
	Price      int       `json:"price"`
	ImageUrl   string    `json:"imageUrl" db:"image_url"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}
