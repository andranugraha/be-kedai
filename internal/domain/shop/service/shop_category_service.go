package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type ShopCategoryService interface {
	GetSellerCategories(userID int, req dto.GetSellerCategoriesRequest) (*commonDto.PaginationResponse, error)
	GetSellerCategoryDetail(userID int, id int) (*dto.ShopCategory, error)
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
