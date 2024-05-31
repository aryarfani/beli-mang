package order

import (
	"beli-mang/internal/entity"
	"time"

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

type CreateOrderRequest struct {
	CalculatedEstimateId uuid.UUID `json:"calculatedEstimateId"`
}

type CreateOrderResponse struct {
	OrderId uuid.UUID `json:"orderId"`
}

// - Param (all optional)
//     - `merchantId` limit the output based on the id
//         - value should be a string
//         - if not exits keep return `200` with empty array
//     - `limit` & `offset` limit the output of the data
//         - default `limit=5&offset=0`
//         - value should be a number
//     - `name` filter based on `merchant’s name` or `item's name`
//         - value should be a string
//         - it should search by wildcard (ex: if search by `name=een` then user with name `kayleen` should appear)
//         - search should be case insensitive
//         - if not exits keep return `200` with empty array
//     - `merchantCategory` filter based on `category`
//         - enum of
//             - `SmallRestaurant`
//             - `MediumRestaurant`
//             - `LargeRestaurant`
//             - `MerchandiseRestaurant`
//             - `BoothKiosk`
//             - `ConvenienceStore`
//         - if not exits / enum is invalid, keep return `200` with empty array

type QueryOrdersRequest struct {
	MerchantId string
	Limit      int
	Offset     int
	Name       string
	Category   string
	UserId     uuid.UUID
}

// [
//   {
//     "orderId": "string",
// 	  "orders": [
// 		  {
// 		    "merchant": {
// 			    "merchantId":"",
// 					"name":"",
// 					"merchantCategory": "",
// 					"imageUrl": "",
// 				  "location": {
// 				    "lat": 1,
// 				    "long": 1
// 				  },
// 				  "createdAt": ""  // should in ISO 8601 format with nanoseconds
// 				},
// 		    "items": [
// 					{
// 						"itemId":"",
// 				    "name": "string",
// 				    "productCategory": ""
// 				    "price": 1,
// 	          "quantity": 1,
// 						"imageUrl": "",
// 					  "createdAt": ""  // should in ISO 8601 format with nanoseconds
// 				  }
// 				]
// 		  }
// 	  ]
//   }
// ]

type QueryOrdersResponse struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserId     uuid.UUID `json:"userId" db:"user_id"`
	MerchantId uuid.UUID `json:"merchantId" db:"merchant_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}
