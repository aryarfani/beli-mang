package merchant

import (
	"beli-mang/internal/entity"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(merchant *entity.Merchant) (merchantId uuid.UUID, err error)
	Query(params QueryMerchantsRequest) (merchants []entity.Merchant, err error)
	QueryNearby(params QueryMerchantsRequest) (QueryResults []QueryMerchantsNearbyResponse, err error)
	Count(params QueryMerchantsRequest) (count int, err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(merchant *entity.Merchant) (merchantId uuid.UUID, err error) {
	query := "INSERT INTO merchants (name, category, image_url, latitude, longitude) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = r.db.QueryRowx(query, merchant.Name, merchant.Category, merchant.ImageUrl, merchant.Latitude, merchant.Longitude).Scan(&merchantId)
	return merchantId, err
}

func (r *repository) Query(params QueryMerchantsRequest) (merchants []entity.Merchant, err error) {
	query := "SELECT * FROM merchants WHERE 1=1"

	if params.MerchantId != "" {
		query += fmt.Sprintf(" AND id = '%s'", params.MerchantId)
	}
	if params.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE '%%%s%%'", params.Name)
	}
	if slices.Contains(entity.MerchantCategories, params.Category) {
		query += fmt.Sprintf(" AND category = '%s'", params.Category)
	}
	if params.Latitude != "" && params.Longitude != "" {
		query += fmt.Sprintf(" ORDER BY ll_to_earth(latitude, longitude) <-> ll_to_earth(%s, %s)", params.Latitude, params.Longitude)
	}
	if params.CreatedAt == "asc" {
		query += " ORDER BY created_at " + params.CreatedAt
	} else {
		query += " ORDER BY created_at DESC"
	}

	limit := 5
	if params.Limit != 0 {
		limit = params.Limit
	}
	offset := 0
	if params.Offset != 0 {
		offset = params.Offset
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	fmt.Println(query)

	err = r.db.Select(&merchants, query)
	return merchants, err
}

type QueryDest struct {
	ID        uuid.UUID `json:"userId"`
	Name      string    `json:"name"`
	Category  string    `json:"merchantCategory"`
	ImageUrl  string    `json:"imageUrl" db:"image_url"`
	Latitude  float64   `json:"lat"`
	Longitude float64   `json:"long"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	Items     string    `db:"items"`
}

func (r *repository) QueryNearby(params QueryMerchantsRequest) (queryResults []QueryMerchantsNearbyResponse, err error) {
	query := `
		SELECT
			merchants.*, jsonb_agg(items) AS items
		FROM
			merchants
			LEFT JOIN items ON items.merchant_id = merchants.id
		WHERE 1=1
	`

	if params.MerchantId != "" {
		query += fmt.Sprintf(" AND merchants.id = '%s'", params.MerchantId)
	}
	if params.Name != "" {
		query += fmt.Sprintf(" AND merchants.name ILIKE '%%%s%%'", params.Name)
	}
	if slices.Contains(entity.MerchantCategories, params.Category) {
		query += fmt.Sprintf(" AND merchants.category = '%s'", params.Category)
	}

	query += "GROUP BY merchants.id"

	if params.Latitude != "" && params.Longitude != "" {
		query += fmt.Sprintf(" ORDER BY ll_to_earth(latitude, longitude) <-> ll_to_earth(%s, %s)", params.Latitude, params.Longitude)
	}
	if params.CreatedAt == "asc" {
		query += " ,created_at " + params.CreatedAt
	} else {
		query += " ,created_at DESC"
	}
	limit := 5
	if params.Limit != 0 {
		limit = params.Limit
	}
	offset := 0
	if params.Offset != 0 {
		offset = params.Offset
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	fmt.Println(query)

	var queryDest []QueryDest
	err = r.db.Select(&queryDest, query)

	for _, queryDest := range queryDest {
		items := []entity.Item{}

		// Check for [null] and set jsonData to an empty array if true
		if string(queryDest.Items) != "[null]" {
			err := json.Unmarshal([]byte(queryDest.Items), &items)
			if err != nil {
				return queryResults, err
			}
		}

		queryResults = append(queryResults, QueryMerchantsNearbyResponse{
			Merchant: QueryMerchantResponse{
				MerchantId:       queryDest.ID.String(),
				Name:             queryDest.Name,
				MerchantCategory: queryDest.Category,
				ImageUrl:         queryDest.ImageUrl,
				CreatedAt:        queryDest.CreatedAt,
				Location: Location{
					Lat:  queryDest.Latitude,
					Long: queryDest.Longitude,
				},
			},
			Items: &items,
		})
	}

	return queryResults, err
}

func (r *repository) Count(params QueryMerchantsRequest) (count int, err error) {
	query := "SELECT COUNT(*) FROM merchants WHERE 1=1"
	if params.MerchantId != "" {
		query += fmt.Sprintf(" AND id = '%s'", params.MerchantId)
	}
	if params.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE '%%%s%%'", params.Name)
	}
	if slices.Contains(entity.MerchantCategories, params.Category) {
		query += fmt.Sprintf(" AND category = '%s'", params.Category)
	}

	err = r.db.Get(&count, query)
	return count, err
}
