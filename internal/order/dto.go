package order

import (
	"beli-mang/internal/entity"
	"time"

	"github.com/google/uuid"
)

type CreateEstimationRequest struct {
	UserLocation UserLocation `json:"userLocation" validate:"required"`
	Orders       []Order      `json:"orders" validate:"dive"`
}

type UserLocation struct {
	Lat  float64 `json:"lat" validate:"required,numeric"`
	Long float64 `json:"long" validate:"required,numeric"`
}

type Order struct {
	MerchantId      string `json:"merchantId" validate:"required"`
	IsStartingPoint bool   `json:"isStartingPoint"`
	Items           []Item `json:"items" validate:"dive"`
}

type Item struct {
	ItemId   string `json:"itemId" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,numeric"`
}

func (req CreateEstimationRequest) ToOrders(userId uuid.UUID) []entity.Order {
	orders := []entity.Order{}

	for _, order := range req.Orders {
		orderItems := []entity.OrderItem{}

		for _, item := range order.Items {
			// only add valid items
			itemUuid, err := uuid.Parse(item.ItemId)
			if err == nil {
				orderItems = append(orderItems, entity.OrderItem{
					ItemId:   itemUuid,
					Quantity: item.Quantity,
				})
			}
		}

		// only add merchants where at least one item valid
		if len(orderItems) > 0 {
			merchantUuid, err := uuid.Parse(order.MerchantId)
			if err == nil {
				orders = append(orders, entity.Order{
					UserId:     userId,
					MerchantId: merchantUuid,
					OrderItems: orderItems,
				})
			}
		}
	}

	return orders
}

type CreateEstimationResponse struct {
	TotalPrice                 int       `json:"totalPrice"`
	EstimatedDeliveryTimeInMin int       `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateId       uuid.UUID `json:"calculatedEstimateId"`
}

type CreateOrderRequest struct {
	CalculatedEstimateId uuid.UUID `json:"calculatedEstimateId" validate:"required,uuid4"`
}

type CreateOrderResponse struct {
	OrderId uuid.UUID `json:"orderId"`
}

type QueryOrdersRequest struct {
	MerchantId       string
	Limit            int
	Offset           int
	Name             string
	MerchantCategory string
	UserId           uuid.UUID
}

type QueryOrdersResponse struct {
	OrderId uuid.UUID   `json:"order_id" db:"order_id"`
	Orders  []OrderResp `json:"orders" db:"orders"`
}

type OrderResp struct {
	MerchantResps MerchantResp `json:"merchant"`
	ItemResps     []ItemResp   `json:"items"`
}

type ItemResp struct {
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	ItemID    string    `json:"itemId"`
	Category  string    `json:"productCategory"`
	Quantity  int       `json:"quantity"`
	ImageURL  string    `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

type MerchantResp struct {
	ID        string    `json:"merchantId"`
	Name      string    `json:"name"`
	Category  string    `json:"merchantCategory"`
	Latitude  float64   `json:"latitude"`
	ImageURL  string    `json:"imageUrl"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"createdAt"`
	Location  Location  `json:"location"`
}
