package merchant

import "beli-mang/internal/entity"

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

type QueryMerchantsResponse struct {
	MerchantId       string   `json:"merchantId"`
	Name             string   `json:"name"`
	MerchantCategory string   `json:"merchantCategory"`
	ImageUrl         string   `json:"imageUrl"`
	Location         Location `json:"location"`
	CreatedAt        string   `json:"createdAt"`
}

func ToQueryMerchantsResponse(merchant *entity.Merchant) *QueryMerchantsResponse {
	return &QueryMerchantsResponse{
		MerchantId:       merchant.ID.String(),
		Name:             merchant.Name,
		MerchantCategory: merchant.Category,
		ImageUrl:         merchant.ImageUrl,
		CreatedAt:        merchant.CreatedAt.Format("2006-01-02 15:04:05"),
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
