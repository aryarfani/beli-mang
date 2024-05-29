package item

type CreateItemRequest struct {
	Name            string `json:"name" validate:"required,min=2,max=30"`
	ProductCategory string `json:"productCategory" validate:"required,oneof=Beverage Food Snack Condiments Additions"`
	Price           int    `json:"price" validate:"required,min=1"`
	ImageUrl        string `json:"imageUrl" validate:"required,xImageUrl"`
	MerchantId      string
}

// - Param (all optional)
//     - `itemId` limit the output based on the id
//         - value should be a string
//         - if not exits keep return `200` with empty array
//     - `limit` & `offset` limit the output of the data
//         - default `limit=5&offset=0`
//         - value should be a number
//     - `name` filter based on name
//         - value should be a string
//         - it should search by wildcard (ex: if search by `name=een` then user with name `kayleen` should appear)
//         - search should be case insensitive
//         - if not exits keep return `200` with empty array
//     - `productCategory` filter based on `category`
//         - enum of
//             - `Beverage`
//             - `Food`
//             - `Snack`
//             - `Condiments`
//             - `Additions`
//         - if not exits / enum is invalid, keep return `200` with empty array
//     - `createdAt` sort by created time
//         - enum of
//             - `asc` sort result from oldest first
//             - `desc` sort result from newest first
//         - if value is wrong, just ignore the param

type QueryItemsRequest struct {
	ItemID     string
	Limit      int
	Offset     int
	Name       string
	Category   string
	CreatedAt  string
	MerchantId string
}
