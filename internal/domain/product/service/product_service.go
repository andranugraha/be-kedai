package service

import (
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/product/repository"
)

type ProductService interface {
	GetByCodeFull(code string) (*model.Product, error)
}

type productServiceImpl struct {
	repository repository.ProductRepository
}

type ProductSConfig struct {
	Repository repository.ProductRepository
}

func NewProductService(cfg *ProductSConfig) ProductService {
	return &productServiceImpl{
		repository: cfg.Repository,
	}
}

func (s *productServiceImpl) GetByCodeFull(code string) (*model.Product, error) {
	return s.repository.GetByCodeFull(code)
}
