package item

type CreateItemRequest struct {
	Name            string `json:"name" validate:"required,min=2,max=30"`
	ProductCategory string `json:"productCategory" validate:"required,oneof=Beverage Food Snack Condiments Additions"`
	Price           int    `json:"price" validate:"required,min=1"`
	ImageUrl        string `json:"imageUrl" validate:"required,xImageUrl"`
}
