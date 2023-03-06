package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/product/repository"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"strings"
)

type ProductService interface {
	GetByID(id int) (*model.Product, error)
	GetByCode(code string) (*dto.ProductDetail, error)
	GetRecommendationByCategory(productId int, categoryId int) ([]*dto.ProductResponse, error)
	ProductSearchFiltering(req dto.ProductSearchFilterRequest) (*commonDto.PaginationResponse, error)
}

type productServiceImpl struct {
	productRepository  repository.ProductRepository
	shopService        service.ShopService
	shopVoucherService service.ShopVoucherService
	courierService     service.CourierService
}

type ProductSConfig struct {
	ProductRepository  repository.ProductRepository
	ShopService        service.ShopService
	ShopVoucherService service.ShopVoucherService
	CourierService     service.CourierService
}

func NewProductService(cfg *ProductSConfig) ProductService {
	return &productServiceImpl{
		productRepository:  cfg.ProductRepository,
		shopVoucherService: cfg.ShopVoucherService,
		courierService:     cfg.CourierService,
		shopService:        cfg.ShopService,
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

func (s *productServiceImpl) ProductSearchFiltering(req dto.ProductSearchFilterRequest) (*commonDto.PaginationResponse, error) {
	validateKeyword := strings.Trim(req.Keyword, " ")
	var shopId int
	if validateKeyword == "" {
		return &commonDto.PaginationResponse{
			Data:       []*dto.ProductResponse{},
			Limit:      req.Limit,
			Page:       req.Page,
			TotalRows:  0,
			TotalPages: 0,
		}, nil
	}

	if req.Shop != "" {
		shop, err := s.shopService.FindShopBySlug(req.Shop)
		if err != nil {
			return nil, err
		}
		shopId = shop.ID
	}

	res, rows, pages, err := s.productRepository.ProductSearchFiltering(req, shopId)
	if err != nil {
		return nil, err
	}

	response := &commonDto.PaginationResponse{
		Data:       res,
		Limit:      req.Limit,
		Page:       req.Page,
		TotalRows:  rows,
		TotalPages: pages,
	}

	return response, nil
}
