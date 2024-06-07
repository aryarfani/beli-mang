package merchant

import (
	"beli-mang/internal/entity"

	"github.com/google/uuid"
)

type Service interface {
	Create(req CreateMerchantRequest) (userId uuid.UUID, err error)
	Query(params QueryMerchantsRequest) (resp PaginatedQueryMerchantsResponse, err error)
	QueryNearby(params QueryMerchantsRequest) (resp PaginatedQueryMerchantsNearbyResponse, err error)
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

func (s *service) Query(params QueryMerchantsRequest) (resp PaginatedQueryMerchantsResponse, err error) {
	merchants, err := s.repo.Query(params)
	if err != nil {
		return resp, err
	}

	merchantsResp := []QueryMerchantResponse{}
	for _, merchant := range merchants {
		merchantsResp = append(merchantsResp, *ToQueryMerchantResponse(&merchant))
	}

	total, err := s.repo.Count(params)
	if err != nil {
		return resp, err
	}

	resp.Data = merchantsResp
	resp.Meta = *entity.NewPaginationMeta(params.Limit, params.Offset, total)

	return resp, nil
}

func (s *service) QueryNearby(params QueryMerchantsRequest) (resp PaginatedQueryMerchantsNearbyResponse, err error) {
	merchants, err := s.repo.QueryNearby(params)
	if err != nil {
		return resp, err
	}

	total, err := s.repo.Count(params)
	if err != nil {
		return resp, err
	}

	resp.Data = merchants
	resp.Meta = *entity.NewPaginationMeta(params.Limit, params.Offset, total)

	return resp, nil
}
