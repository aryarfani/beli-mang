package item

import (
	"beli-mang/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(item *CreateItemRequest) (itemId uuid.UUID, err error)
	Query() (items []entity.Item, err error)
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
	query := "INSERT INTO items (name, category, price, image_url) VALUES ($1, $2, $3, $4) RETURNING id"
	err = r.db.QueryRowx(query, item.Name, item.ProductCategory, item.Price, item.ImageUrl).Scan(&itemId)
	return itemId, err
}

func (r *repository) Query() (items []entity.Item, err error) {
	query := "SELECT * FROM items"
	err = r.db.Select(&items, query)
	return items, err
}
