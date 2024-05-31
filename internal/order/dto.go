package order

import (
	"beli-mang/internal/entity"

	"github.com/google/uuid"
)

type CreateEstimationRequest struct {
	UserLocation UserLocation `json:"userLocation"`
	Orders       []Order      `json:"orders"`
}

type UserLocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Order struct {
	MerchantId      uuid.UUID `json:"merchantId"`
	IsStartingPoint bool      `json:"isStartingPoint"`
	Items           []Item    `json:"items"`
}

type Item struct {
	ItemId   string `json:"itemId"`
	Quantity int    `json:"quantity"`
}

func (req CreateEstimationRequest) ToOrders(userId uuid.UUID) []entity.Order {
	orders := []entity.Order{}

	for _, order := range req.Orders {
		orderItems := []entity.OrderItem{}

		for _, item := range order.Items {
			orderItems = append(orderItems, entity.OrderItem{
				ItemId:   uuid.MustParse(item.ItemId),
				Quantity: item.Quantity,
			})
		}

		orders = append(orders, entity.Order{
			UserId:     userId,
			MerchantId: order.MerchantId,
			OrderItems: orderItems,
		})
	}

	return orders
}

type CreateEstimationResponse struct {
	TotalPrice                 int       `json:"totalPrice"`
	EstimatedDeliveryTimeInMin int       `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateId       uuid.UUID `json:"calculatedEstimateId"`
}
