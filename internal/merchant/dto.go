package merchant

import (
	"beli-mang/internal/entity"
	"time"
)

type CreateMerchantRequest struct {
	Name             string   `json:"name" validate:"required,min=2,max=30"`
	MerchantCategory string   `json:"merchantCategory" validate:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`
	ImageUrl         string   `json:"imageUrl" validate:"required,xImageUrl"`
	Location         Location `json:"location" validate:"required"`
}

type Location struct {
	Lat  float64 `json:"lat" validate:"required,numeric"`
	Long float64 `json:"long" validate:"required,numeric"`
}

func (req *CreateMerchantRequest) ToMerchant() *entity.Merchant {
	return &entity.Merchant{
		Name:      req.Name,
		Category:  req.MerchantCategory,
		ImageUrl:  req.ImageUrl,
		Latitude:  req.Location.Lat,
		Longitude: req.Location.Long,
	}
}

type PaginatedQueryMerchantsResponse struct {
	Data []QueryMerchantResponse `json:"data"`
	Meta entity.PaginationMeta   `json:"meta"`
}

type PaginatedQueryMerchantsNearbyResponse struct {
	Data []QueryMerchantsNearbyResponse `json:"data"`
	Meta entity.PaginationMeta          `json:"meta"`
}

type QueryMerchantsNearbyResponse struct {
	Merchant QueryMerchantResponse `json:"merchant"`
	Items    *[]entity.Item         `json:"item"`
}

type QueryMerchantResponse struct {
	MerchantId       string    `json:"merchantId"`
	Name             string    `json:"name"`
	MerchantCategory string    `json:"merchantCategory"`
	ImageUrl         string    `json:"imageUrl"`
	Location         Location  `json:"location"`
	CreatedAt        time.Time `json:"createdAt"`
}

func ToQueryMerchantResponse(merchant *entity.Merchant) *QueryMerchantResponse {
	return &QueryMerchantResponse{
		MerchantId:       merchant.ID.String(),
		Name:             merchant.Name,
		MerchantCategory: merchant.Category,
		ImageUrl:         merchant.ImageUrl,
		CreatedAt:        merchant.CreatedAt,
		Location: Location{
			Lat:  merchant.Latitude,
			Long: merchant.Longitude,
		},
	}
}

type QueryMerchantsRequest struct {
	MerchantId string `query:"merchantId"`
	Limit      int    `query:"limit"`
	Offset     int    `query:"offset"`
	Name       string `query:"name"`
	Category   string `query:"merchantCategory"`
	CreatedAt  string `query:"createdAt"`
	Latitude   string
	Longitude  string
}
