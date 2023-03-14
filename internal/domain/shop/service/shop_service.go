package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
	"strings"
)

type ShopService interface {
	CreateShop(userID int, request *dto.CreateShopRequest) (*model.Shop, error)
	FindShopById(id int) (*model.Shop, error)
	FindShopByUserId(userId int) (*model.Shop, error)
	FindShopBySlug(slug string) (*model.Shop, error)
	FindShopByKeyword(req dto.FindShopRequest) (*commonDto.PaginationResponse, error)
	GetShopFinanceOverview(userId int) (*dto.ShopFinanceOverviewResponse, error)
}

type shopServiceImpl struct {
	shopRepository        repository.ShopRepository
	courierServiceService CourierServiceService
}

type ShopSConfig struct {
	ShopRepository        repository.ShopRepository
	CourierServiceService CourierServiceService
}

func NewShopService(cfg *ShopSConfig) ShopService {
	return &shopServiceImpl{
		shopRepository:        cfg.ShopRepository,
		courierServiceService: cfg.CourierServiceService,
	}
}

func (s *shopServiceImpl) FindShopById(id int) (*model.Shop, error) {
	return s.shopRepository.FindShopById(id)
}

func (s *shopServiceImpl) FindShopByUserId(userId int) (*model.Shop, error) {
	return s.shopRepository.FindShopByUserId(userId)
}

func (s *shopServiceImpl) FindShopBySlug(slug string) (*model.Shop, error) {
	shop, err := s.shopRepository.FindShopBySlug(slug)
	if err != nil {
		return nil, err
	}

	return shop, nil
}

func (s *shopServiceImpl) FindShopByKeyword(req dto.FindShopRequest) (*commonDto.PaginationResponse, error) {
	validateKeyword := strings.Trim(req.Keyword, " ")
	if validateKeyword == "" {
		return &commonDto.PaginationResponse{
			Data:       []*dto.FindShopResponse{},
			TotalRows:  0,
			TotalPages: 0,
			Limit:      req.Limit,
			Page:       req.Page,
		}, nil
	}

	res, rows, pages, err := s.shopRepository.FindShopByKeyword(req)
	if err != nil {
		return nil, err
	}

	return &commonDto.PaginationResponse{
		Data:       res,
		TotalRows:  rows,
		TotalPages: pages,
		Limit:      req.Limit,
		Page:       req.Page,
	}, nil
}

func (s *shopServiceImpl) GetShopFinanceOverview(userId int) (*dto.ShopFinanceOverviewResponse, error) {
	shop, err := s.shopRepository.FindShopByUserId(userId)
	if err != nil {
		return nil, err
	}

	shopFinanceOverview, err := s.shopRepository.GetShopFinanceOverview(shop.ID)
	if err != nil {
		return nil, err
	}

	return shopFinanceOverview, nil
}

func (s *shopServiceImpl) CreateShop(userID int, request *dto.CreateShopRequest) (*model.Shop, error) {
	var shop model.Shop

	return &shop, nil
}
