package item

import (
	"beli-mang/internal/entity"

	"github.com/google/uuid"
)

type Service interface {
	Create(req CreateItemRequest) (userId uuid.UUID, err error)
	Query() (items []entity.Item, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(req CreateItemRequest) (userId uuid.UUID, err error) {
	userId, err = s.repo.Create(&req)
	if err != nil {
		return userId, err
	}

	return userId, nil
}

func (s *service) Query() (items []entity.Item, err error) {
	items, err = s.repo.Query()
	if err != nil {
		return items, err
	}

	return items, nil
}
