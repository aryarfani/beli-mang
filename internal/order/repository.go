package order

import (
	"beli-mang/internal/entity"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(orders []entity.Order, estimation entity.Estimation) (estimationId uuid.UUID, err error)
	GetItems(itemIds []uuid.UUID) (items []entity.Item, err error)
	GetMerchants(itemIds []uuid.UUID) (merchants []entity.Merchant, err error)
	CreateOrder(req CreateOrderRequest) (orderId uuid.UUID, err error)
	Query(params QueryOrdersRequest) (resp []QueryOrdersResponse, err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(orders []entity.Order, estimation entity.Estimation) (estimationId uuid.UUID, err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return estimationId, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// insert estimation
	query := "INSERT INTO estimations (total_price, estimated_delivery_time_in_minutes) VALUES ($1, $2) RETURNING id"
	err = tx.QueryRowx(query, estimation.TotalPrice, estimation.EstimatedDeliveryTimeInMinutes).Scan(&estimationId)
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

func (r *repository) GetMerchants(itemIds []uuid.UUID) (merchants []entity.Merchant, err error) {
	query, args, err := sqlx.In("SELECT * FROM merchants WHERE id IN (?)", itemIds)
	if err != nil {
		return merchants, err
	}
	query = r.db.Rebind(query)
	err = r.db.Select(&merchants, query, args...)
	return merchants, err
}

func (r *repository) CreateOrder(req CreateOrderRequest) (orderId uuid.UUID, err error) {
	query := "UPDATE orders SET estimation_id = null WHERE estimation_id = $1 RETURNING id"
	err = r.db.QueryRowx(query, req.CalculatedEstimateId).Scan(&orderId)
	return orderId, err
}

func (r *repository) Query(params QueryOrdersRequest) (resp []QueryOrdersResponse, err error) {
	query := "SELECT * FROM orders "
	query += fmt.Sprintf("WHERE user_id = '%s'", params.UserId)

	if params.MerchantId != "" {
		query += fmt.Sprintf(" AND merchant_id = '%s'", params.MerchantId)
	}
	// if params.Name != "" {
	// 	query += fmt.Sprintf(" AND name ILIKE '%%%s%%'", params.Name)
	// }
	// if slices.Contains(entity.MerchantCategories, params.Category) {
	// 	query += fmt.Sprintf(" AND category = '%s'", params.Category)
	// }

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

	err = r.db.Select(&resp, query)
	return resp, err
}
