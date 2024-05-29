package order

import "beli-mang/internal/entity"

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

	resp.CalculatedEstimateId = estimationId

	return resp, nil
}
