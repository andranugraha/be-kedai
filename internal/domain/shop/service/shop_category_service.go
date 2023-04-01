package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type ShopCategoryService interface {
	GetSellerCategories(userID int, req dto.GetSellerCategoriesRequest) (*commonDto.PaginationResponse, error)
	GetSellerCategoryDetail(userID int, id int) (*dto.ShopCategory, error)
	CreateSellerCategory(userID int, req dto.CreateSellerCategoryRequest) (*dto.CreateSellerCategoryResponse, error)
	UpdateSellerCategory(userID int, id int, req dto.UpdateSellerCategoryRequest) (*dto.CreateSellerCategoryResponse, error)
}

type shopCategoryServiceImpl struct {
	shopService      ShopService
	shopCategoryRepo repository.ShopCategoryRepository
}

type ShopCategorySConfig struct {
	ShopService      ShopService
	ShopCategoryRepo repository.ShopCategoryRepository
}

func NewShopCategoryService(cfg *ShopCategorySConfig) ShopCategoryService {
	return &shopCategoryServiceImpl{
		shopService:      cfg.ShopService,
		shopCategoryRepo: cfg.ShopCategoryRepo,
	}
}

func (s *shopCategoryServiceImpl) GetSellerCategories(userID int, req dto.GetSellerCategoriesRequest) (*commonDto.PaginationResponse, error) {
	shop, err := s.shopService.FindShopById(userID)
	if err != nil {
		return nil, err
	}

	shopCategories, totalRows, totalPages, err := s.shopCategoryRepo.GetByShopID(shop.ID, req)
	if err != nil {
		return nil, err
	}

	return &commonDto.PaginationResponse{
		Data:       shopCategories,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Page:       req.Page,
		Limit:      req.Limit,
	}, nil
}

func (s *shopCategoryServiceImpl) GetSellerCategoryDetail(userID int, id int) (*dto.ShopCategory, error) {
	shop, err := s.shopService.FindShopById(userID)
	if err != nil {
		return nil, err
	}

	shopCategory, err := s.shopCategoryRepo.GetByIDAndShopID(id, shop.ID)
	if err != nil {
		return nil, err
	}

	return shopCategory, nil
}

func (s *shopCategoryServiceImpl) CreateSellerCategory(userID int, req dto.CreateSellerCategoryRequest) (*dto.CreateSellerCategoryResponse, error) {
	shop, err := s.shopService.FindShopById(userID)
	if err != nil {
		return nil, err
	}

	shopCategory := req.ComposeModel(shop.ID)

	err = s.shopCategoryRepo.Create(shopCategory)
	if err != nil {
		return nil, err
	}

	return &dto.CreateSellerCategoryResponse{
		ID: shopCategory.ID,
	}, nil
}

func (s *shopCategoryServiceImpl) UpdateSellerCategory(userID int, id int, req dto.UpdateSellerCategoryRequest) (*dto.CreateSellerCategoryResponse, error) {
	shop, err := s.shopService.FindShopById(userID)
	if err != nil {
		return nil, err
	}

	category, err := s.shopCategoryRepo.GetCategoryByIDAndShopID(id, shop.ID)
	if err != nil {
		return nil, err
	}

	category.Products = func() []*model.ShopCategoryProduct {
		var products []*model.ShopCategoryProduct
		for _, productId := range req.ProductIDs {
			products = append(products, &model.ShopCategoryProduct{
				ProductId:      productId,
				ShopCategoryId: category.ID,
			})
		}
		return products
	}()

	err = s.shopCategoryRepo.Update(category)
	if err != nil {
		return nil, err
	}

	return &dto.CreateSellerCategoryResponse{
		ID: category.ID,
	}, nil
}
