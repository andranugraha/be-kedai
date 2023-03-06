package service

import (
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/product/repository"
)

type SkuService interface {
	GetByID(id int) (*model.Sku, error)
	GetSKUByVariantIDs(request *dto.GetSKURequest) (*model.Sku, error)
}

type skuServiceImpl struct {
	skuRepository repository.SkuRepository
}

type SkuSConfig struct {
	SkuRepository repository.SkuRepository
}

func NewSkuService(cfg *SkuSConfig) SkuService {
	return &skuServiceImpl{
		skuRepository: cfg.SkuRepository,
	}
}

func (s *skuServiceImpl) GetByID(id int) (*model.Sku, error) {
	return s.skuRepository.GetByID(id)
}

func (s *skuServiceImpl) GetSKUByVariantIDs(request *dto.GetSKURequest) (*model.Sku, error) {
	variantIDs, err := request.ToIntList()
	if err != nil {
		return nil, err
	}

	return s.skuRepository.GetByVariantIDs(variantIDs)
}
