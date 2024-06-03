package item

import (
	"beli-mang/internal/entity"

	"github.com/google/uuid"
)

type Service interface {
	Create(req CreateItemRequest) (userId uuid.UUID, err error)
	Query(params QueryItemsRequest) (resp PaginatedQueryItemsResponse, err error)
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

func (s *service) Query(params QueryItemsRequest) (resp PaginatedQueryItemsResponse, err error) {
	items, err := s.repo.Query(params)
	if err != nil {
		return resp, err
	}

	total, err := s.repo.Count(params)
	if err != nil {
		return resp, err
	}

	resp.Data = items
	resp.Meta = *entity.NewPaginationMeta(params.Limit, params.Offset, total)

	return resp, nil
}
