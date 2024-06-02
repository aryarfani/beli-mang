package order

import (
	"beli-mang/internal/entity"
	"encoding/json"
	"fmt"
	"log"
	"slices"

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
		log.Println(order.UserId)
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

type QueryResult struct {
	OrderId uuid.UUID `db:"order_id"`
	Orders  string    `db:"orders"`
}

func (r *repository) Query(params QueryOrdersRequest) (resp []QueryOrdersResponse, err error) {
	query := `
		SELECT
			e.id order_id,
			jsonb_agg(
				jsonb_build_object(
					'merchant', (m.*),
					'items', (
								SELECT
										jsonb_agg(
											json_build_object(
												'item_id', i.id,
												'name', i.name,
												'category', i.category,
												'price', i.price,
												'quantity', oi.quantity,
												'image_url', i.image_url,
												'created_at', i.created_at
											)
										)
								FROM
									order_items oi
									JOIN items i ON oi.item_id = i.id
									AND oi.order_id = o.id
								)
							)
			) orders
		FROM
			estimations e
			LEFT JOIN orders o ON e.id = o.estimation_id
			LEFT JOIN merchants m ON o.merchant_id = m.id
			LEFT JOIN order_items oi ON oi.order_id = o.id
			LEFT JOIN items i ON oi.item_id = i.id
		WHERE 1=1
	`
	query += fmt.Sprintf("AND o.user_id = '%s'", params.UserId)

	if params.MerchantId != "" {
		query += fmt.Sprintf(" AND o.merchant_id = '%s'", params.MerchantId)
	}
	if params.Name != "" {
		query += fmt.Sprintf(` AND (m.name ILIKE '%%%s%%' OR i.name ILIKE '%%%s%%')`, params.Name, params.Name)
	}
	if slices.Contains(entity.MerchantCategories, params.MerchantCategory) {
		query += fmt.Sprintf(" AND m.category = '%s'", params.MerchantCategory)
	}

	limit := 5
	if params.Limit != 0 {
		limit = params.Limit
	}
	offset := 0
	if params.Offset != 0 {
		offset = params.Offset
	}

	query += fmt.Sprintf(" GROUP BY e.id LIMIT %d OFFSET %d", limit, offset)
	fmt.Println(query)

	var queryResults []QueryResult
	err = r.db.Select(&queryResults, query)

	for _, queryResult := range queryResults {
		var order []OrderResp
		err := json.Unmarshal([]byte(queryResult.Orders), &order)
		if err != nil {
			return resp, err
		}

		resp = append(resp, QueryOrdersResponse{
			OrderId: queryResult.OrderId,
			Orders:  order,
		})
	}
	return resp, err
}
