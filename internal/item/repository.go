package item

import (
	"beli-mang/internal/entity"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(item *CreateItemRequest) (itemId uuid.UUID, err error)
	Query(params QueryItemsRequest) (items []entity.Item, err error)
	Count(params QueryItemsRequest) (count int, err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(item *CreateItemRequest) (itemId uuid.UUID, err error) {
	query := "INSERT INTO items (merchant_id, name, category, price, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = r.db.QueryRowx(query, item.MerchantId, item.Name, item.ProductCategory, item.Price, item.ImageUrl).Scan(&itemId)
	return itemId, err
}

func (r *repository) Query(params QueryItemsRequest) (items []entity.Item, err error) {
	query := "SELECT * FROM items"
	query += fmt.Sprintf(" WHERE merchant_id = '%s'", params.MerchantId)

	if params.ItemID != "" {
		query += fmt.Sprintf(" AND id = '%s'", params.ItemID)
	}
	if params.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE '%%%s%%'", params.Name)
	}
	if slices.Contains(entity.ItemCategories, params.Category) {
		query += fmt.Sprintf(" AND category = '%s'", params.Category)
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

	err = r.db.Select(&items, query)
	return items, err
}

func (r *repository) Count(params QueryItemsRequest) (count int, err error) {
	query := "SELECT COUNT(*) FROM items"
	query += fmt.Sprintf(" WHERE merchant_id = '%s'", params.MerchantId)

	if params.ItemID != "" {
		query += fmt.Sprintf(" AND id = '%s'", params.ItemID)
	}
	if params.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE '%%%s%%'", params.Name)
	}
	if slices.Contains(entity.ItemCategories, params.Category) {
		query += fmt.Sprintf(" AND category = '%s'", params.Category)
	}

	err = r.db.Get(&count, query)
	return count, err
}
