package service

import (
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/product/repository"
)

type SkuService interface {
	GetByID(id int) (*model.Sku, error)
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
