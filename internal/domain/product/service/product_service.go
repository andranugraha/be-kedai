package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/product/repository"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	productUtils "kedai/backend/be-kedai/internal/utils/product"
	"strings"
)

type ProductService interface {
	GetByID(id int) (*model.Product, error)
	GetActiveByID(id int) (*model.Product, error)
	GetByCode(code string) (*dto.ProductDetail, error)
	GetProductsByShopSlug(slug string, request *dto.ShopProductFilterRequest) (*commonDto.PaginationResponse, error)
	GetRecommendationByCategory(productId int, categoryId int) ([]*dto.ProductResponse, error)
	ProductSearchFiltering(req dto.ProductSearchFilterRequest) (*commonDto.PaginationResponse, error)
	GetSellerProducts(userID int, req *dto.SellerProductFilterRequest) (*commonDto.PaginationResponse, error)
	SearchAutocomplete(req dto.ProductSearchAutocomplete) ([]*dto.ProductResponse, error)
	GetSellerProductByCode(userID int, productCode string) (*dto.SellerProductDetail, error)
	AddViewCount(id int) error
	UpdateProductActivation(userID int, code string, request *dto.UpdateProductActivationRequest) error
	CreateProduct(userID int, request *dto.CreateProductRequest) (*model.Product, error)
	GetRecommendedProducts(req *dto.GetRecommendedProductRequest) (*commonDto.PaginationResponse, error)
	UpdateProduct(userID int, code string, request *dto.CreateProductRequest) (*model.Product, error)
}

type productServiceImpl struct {
	productRepository     repository.ProductRepository
	shopService           service.ShopService
	shopVoucherService    service.ShopVoucherService
	courierService        service.CourierService
	courierServiceService service.CourierServiceService
	categoryService       CategoryService
	DiscussionService     DiscussionService
}

type ProductSConfig struct {
	ProductRepository     repository.ProductRepository
	ShopService           service.ShopService
	ShopVoucherService    service.ShopVoucherService
	CourierService        service.CourierService
	CourierServiceService service.CourierServiceService
	CategoryService       CategoryService
	DiscussionService     DiscussionService
}

func NewProductService(cfg *ProductSConfig) ProductService {
	return &productServiceImpl{
		productRepository:     cfg.ProductRepository,
		shopVoucherService:    cfg.ShopVoucherService,
		courierService:        cfg.CourierService,
		shopService:           cfg.ShopService,
		categoryService:       cfg.CategoryService,
		courierServiceService: cfg.CourierServiceService,
		DiscussionService:     cfg.DiscussionService,
	}
}

func (s *productServiceImpl) GetByID(id int) (*model.Product, error) {
	return s.productRepository.GetByID(id)
}

func (s *productServiceImpl) GetActiveByID(id int) (*model.Product, error) {
	return s.productRepository.GetActiveByID(id)
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

	couriers, err := s.courierService.GetCouriersByProductID(productDetail.ID)
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
	if validateKeyword == "" && req.CategoryId == 0 {
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

	var categoryIDs []int
	if req.CategoryId > 0 {
		cIDs, err := s.categoryService.GetCategoryIDLineAgesFromTop(req.CategoryId)
		if err != nil {
			return nil, err
		}

		categoryIDs = cIDs
	}

	res, rows, pages, err := s.productRepository.ProductSearchFiltering(req, shopId, categoryIDs)
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

func (s *productServiceImpl) GetProductsByShopSlug(slug string, request *dto.ShopProductFilterRequest) (*commonDto.PaginationResponse, error) {
	shop, err := s.shopService.FindShopBySlug(slug)

	if err != nil {
		return nil, err
	}

	products, totalRows, totalPages, err := s.productRepository.GetByShopID(shop.ID, request)

	if err != nil {
		return nil, err
	}

	response := commonDto.PaginationResponse{
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Page:       request.Page,
		Limit:      request.Limit,
		Data:       products,
	}

	return &response, nil
}

func (s *productServiceImpl) GetSellerProducts(userID int, req *dto.SellerProductFilterRequest) (*commonDto.PaginationResponse, error) {
	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return nil, err
	}

	products, totalRows, totalPages, err := s.productRepository.GetBySellerID(shop.ID, req)
	if err != nil {
		return nil, err
	}

	return &commonDto.PaginationResponse{
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Page:       req.Page,
		Limit:      req.Limit,
		Data:       products,
	}, nil
}

func (s *productServiceImpl) SearchAutocomplete(req dto.ProductSearchAutocomplete) ([]*dto.ProductResponse, error) {
	return s.productRepository.SearchAutocomplete(req)
}

func (s *productServiceImpl) GetSellerProductByCode(userID int, productCode string) (*dto.SellerProductDetail, error) {
	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepository.GetSellerProductByCode(shop.ID, productCode)
	if err != nil {
		return nil, err
	}

	categories, err := s.categoryService.GetCategoryLineAgesFromBottom(product.CategoryID)
	if err != nil {
		return nil, err
	}

	couriers, err := s.courierService.GetCouriersByProductID(product.ID)
	if err != nil {
		return nil, err
	}

	res := dto.SellerProductDetail{
		Product:    *product,
		Categories: categories,
		Couriers:   couriers,
	}

	return &res, nil
}

func (s *productServiceImpl) AddViewCount(id int) error {
	return s.productRepository.AddViewCount(id)
}

func (s *productServiceImpl) UpdateProductActivation(userID int, code string, request *dto.UpdateProductActivationRequest) error {
	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return err
	}

	return s.productRepository.UpdateActivation(shop.ID, code, *request.IsActive)
}

func (s *productServiceImpl) CreateProduct(userID int, request *dto.CreateProductRequest) (*model.Product, error) {
	if isProductNameValid := productUtils.ValidateProductName(request.Name); !isProductNameValid {
		return nil, commonErr.ErrInvalidProductNamePattern
	}

	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return nil, err
	}

	couriers, err := s.courierServiceService.GetCourierServicesByCourierIDs(request.CourierIDs)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepository.Create(shop.ID, request, couriers)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productServiceImpl) GetRecommendedProducts(request *dto.GetRecommendedProductRequest) (*commonDto.PaginationResponse, error) {

	var resp commonDto.PaginationResponse

	recommendedProducts, totalRows, totalPages, err := s.productRepository.GetRecommended(request)
	if err != nil {
		return nil, err
	}

	resp = commonDto.PaginationResponse{
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Page:       request.Page,
		Limit:      request.Limit,
		Data:       recommendedProducts,
	}

	return &resp, nil
}

func (s *productServiceImpl) UpdateProduct(userID int, code string, request *dto.CreateProductRequest) (*model.Product, error) {
	if isProductNameValid := productUtils.ValidateProductName(request.Name); !isProductNameValid {
		return nil, commonErr.ErrInvalidProductNamePattern
	}

	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return nil, err
	}

	couriers, err := s.courierServiceService.GetCourierServicesByCourierIDs(request.CourierIDs)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepository.Update(shop.ID, code, request, couriers)
	if err != nil {
		return nil, err
	}

	return product, nil
}
