package service

import (
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/product/repository"
)

type ProductService interface {
	GetByID(id int) (*model.Product, error)
	GetByCode(code string) (*model.Product, error)
	GetRecommendationByCategory(productId int, categoryId int) ([]*dto.ProductResponse, error)
}

type productServiceImpl struct {
	productRepository repository.ProductRepository
}

type ProductSConfig struct {
	ProductRepository repository.ProductRepository
}

func NewProductService(cfg *ProductSConfig) ProductService {
	return &productServiceImpl{
		productRepository: cfg.ProductRepository,
	}
}

func (s *productServiceImpl) GetByID(id int) (*model.Product, error) {
	return s.productRepository.GetByID(id)
}

func (s *productServiceImpl) GetByCode(code string) (*model.Product, error) {
	return s.productRepository.GetByCode(code)
}

func (s *productServiceImpl) GetRecommendationByCategory(productId int, categoryId int) ([]*dto.ProductResponse, error) {
	return s.productRepository.GetRecommendationByCategory(productId, categoryId)
}
