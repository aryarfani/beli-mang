package order

import (
	"beli-mang/internal/entity"

	"github.com/google/uuid"
)

// {
//   "userLocation": {
//     "lat": 1, // not null | float
//     "long": 1  // not null | float
//   },
//   "orders": [
//     {
//       "merchantId": "string", // not null
//       "isStartingPoint" : true
//        ⬆️ not null | there's should be one isStartingPoint == true in orders array
//        | if none are true, or true > 1 items, it's not valid
//       "items": [
//         {
//           "itemId": "string", // not null
//           "quantity": 1 // not null
//         }
//       ]
//     }
//   ]
// }

type CreateEstimationRequest struct {
	UserLocation UserLocation `json:"userLocation"`
	Orders       []Order      `json:"orders"`
}

type UserLocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Order struct {
	MerchantId      string `json:"merchantId"`
	IsStartingPoint bool   `json:"isStartingPoint"`
	Items           []Item `json:"items"`
}

type Item struct {
	ItemId   string `json:"itemId"`
	Quantity int    `json:"quantity"`
}

func (req CreateEstimationRequest) ToOrder(userId uuid.UUID) []entity.Order {
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

// {
// 	"totalPrice": 1,
// 	"estimatedDeliveryTimeInMinutes": 1,
// 	"calculatedEstimateId": "" // save the calculation in the system
// }

type CreateEstimationResponse struct {
	TotalPrice                 float64   `json:"totalPrice"`
	EstimatedDeliveryTimeInMin int       `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateId       uuid.UUID `json:"calculatedEstimateId"`
}
