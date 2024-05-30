package order

import (
	"beli-mang/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(orders []entity.Order) (estimationId uuid.UUID, err error)
	GetItems(itemIds []uuid.UUID) (items []entity.Item, err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(orders []entity.Order) (estimationId uuid.UUID, err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return estimationId, err
	}

	// Rollback the transaction in case of an error
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // if err is nil, Commit; otherwise, Rollback
		}
	}()

	// insert estimation
	query := "INSERT INTO estimations (total_price, estimated_delivery_time_in_minutes) VALUES (10,20) RETURNING id"
	err = tx.QueryRowx(query).Scan(&estimationId)
	if err != nil {
		return estimationId, err
	}

	// insert orders
	for _, order := range orders {
		var orderId string
		query := "INSERT INTO orders (user_id, merchant_id, estimation_id) VALUES ($1, $2, $3) RETURNING id"
		err := tx.QueryRowx(query, order.UserId, order.MerchantId, estimationId).Scan(&orderId)
		if err != nil {
			return estimationId, err
		}

		// insert order_items
		for _, item := range order.OrderItems {
			query := "INSERT INTO order_items (order_id, item_id, quantity) VALUES ($1, $2, $3)"
			_, err = tx.Exec(query, orderId, item.ItemId, item.Quantity)
			if err != nil {
				return estimationId, err
			}
		}
	}

	return estimationId, nil
}

func (r *repository) GetItems(itemIds []uuid.UUID) (items []entity.Item, err error) {
	query, args, err := sqlx.In("SELECT * FROM items WHERE id IN (?)", itemIds)
	if err != nil {
		return items, err
	}
	query = r.db.Rebind(query)
	err = r.db.Select(&items, query, args...)
	return items, err
}
