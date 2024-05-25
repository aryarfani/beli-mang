package merchant

import (
	"beli-mang/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(merchant *entity.Merchant) (merchantId uuid.UUID, err error)
	Query() (merchants []entity.Merchant, err error)
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

func (r *repository) Query() (merchants []entity.Merchant, err error) {
	query := "SELECT * FROM merchants"
	err = r.db.Select(&merchants, query)
	return merchants, err
}
