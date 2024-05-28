package merchant

import (
	"github.com/google/uuid"
)

type Service interface {
	Create(req CreateMerchantRequest) (userId uuid.UUID, err error)
	Query(params QueryMerchantsRequest) (merchants []QueryMerchantsResponse, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(req CreateMerchantRequest) (userId uuid.UUID, err error) {
	merchant := req.ToMerchant()
	userId, err = s.repo.Create(merchant)
	if err != nil {
		return userId, err
	}

	return userId, nil
}

func (s *service) Query(params QueryMerchantsRequest) (queryMerchants []QueryMerchantsResponse, err error) {
	merchants, err := s.repo.Query(params)
	if err != nil {
		return queryMerchants, err
	}

	for _, merchant := range merchants {
		queryMerchants = append(queryMerchants, *ToQueryMerchantsResponse(&merchant))
	}

	if len(queryMerchants) == 0 {
		return []QueryMerchantsResponse{}, nil
	}

	return queryMerchants, nil
}
