package service

import (
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/product/repository"
	"kedai/backend/be-kedai/internal/domain/shop/service"
)

type ProductService interface {
	GetByID(id int) (*model.Product, error)
	GetByCode(code string) (*dto.ProductDetail, error)
	GetRecommendationByCategory(productId int, categoryId int) ([]*dto.ProductResponse, error)
}

type productServiceImpl struct {
	productRepository  repository.ProductRepository
	shopVoucherService service.ShopVoucherService
	courierService     service.CourierService
}

type ProductSConfig struct {
	ProductRepository  repository.ProductRepository
	ShopVoucherService service.ShopVoucherService
	CourierService     service.CourierService
}

func NewProductService(cfg *ProductSConfig) ProductService {
	return &productServiceImpl{
		productRepository:  cfg.ProductRepository,
		shopVoucherService: cfg.ShopVoucherService,
		courierService:     cfg.CourierService,
	}
}

func (s *productServiceImpl) GetByID(id int) (*model.Product, error) {
	return s.productRepository.GetByID(id)
}

func (s *productServiceImpl) GetByCode(code string) (*dto.ProductDetail, error) {
	productDetail, err := s.productRepository.GetByCode(code)
	if err != nil {
		return nil, err
	}

	vouchers, err := s.shopVoucherService.GetShopVoucher(productDetail.Shop.Slug)
	if err == nil {
		productDetail.Vouchers = vouchers
	}

	couriers, err := s.courierService.GetCouriersByShopID(productDetail.ShopID)
	if err == nil {
		productDetail.Couriers = couriers
	}

	return productDetail, nil
}

func (s *productServiceImpl) GetRecommendationByCategory(productId int, categoryId int) ([]*dto.ProductResponse, error) {
	return s.productRepository.GetRecommendationByCategory(productId, categoryId)
}
