package order

import (
	"beli-mang/internal/entity"

	"github.com/google/uuid"
)

type Service interface {
	Create(orders []entity.Order) (resp CreateEstimationResponse, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(orders []entity.Order) (resp CreateEstimationResponse, err error) {
	estimationId, err := s.repo.Create(orders)
	if err != nil {
		return resp, err
	}

	// get item uuids from req
	var itemUuids []uuid.UUID
	for _, order := range orders {
		for _, item := range order.OrderItems {
			itemUuids = append(itemUuids, item.ItemId)
		}
	}

	// get items form uuids
	items, err := s.repo.GetItems(itemUuids)
	if err != nil {
		return resp, err
	}

	// create product map for faster search
	var itemsMap = make(map[uuid.UUID]entity.Item)
	for _, item := range items {
		itemsMap[item.ID] = item
	}

	var totalPrice int
	for _, order := range orders {
		for _, orderItem := range order.OrderItems {
			itemPrice := itemsMap[orderItem.ItemId].Price
			totalPrice += itemPrice * orderItem.Quantity
		}
	}

	resp.CalculatedEstimateId = estimationId
	resp.TotalPrice = totalPrice

	return resp, nil
}
