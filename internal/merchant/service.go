package merchant

import (
	"github.com/google/uuid"
)

type Service interface {
	Create(req CreateMerchantRequest) (userId uuid.UUID, err error)
	Query() (merchants []QueryMerchantsResponse, err error)
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

func (s *service) Query() (queryMerchants []QueryMerchantsResponse, err error) {
	merchants, err := s.repo.Query()
	if err != nil {
		return queryMerchants, err
	}

	for _, merchant := range merchants {
		queryMerchants = append(queryMerchants, *ToQueryMerchantsResponse(&merchant))
	}

	return queryMerchants, nil
}
