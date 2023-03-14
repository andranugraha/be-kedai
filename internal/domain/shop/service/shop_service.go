package service

import (
	"errors"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	commonError "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
	shopUtils "kedai/backend/be-kedai/internal/utils/shop"
	stringUtils "kedai/backend/be-kedai/internal/utils/strings"
	"strings"
	"time"
)

type ShopService interface {
	CreateShop(userID int, request *dto.CreateShopRequest) (*model.Shop, error)
	FindShopById(id int) (*model.Shop, error)
	FindShopByUserId(userId int) (*model.Shop, error)
	FindShopBySlug(slug string) (*model.Shop, error)
	FindShopByKeyword(req dto.FindShopRequest) (*commonDto.PaginationResponse, error)
	GetShopFinanceOverview(userId int) (*dto.ShopFinanceOverviewResponse, error)
	GetShopStats(userId int) (*dto.GetShopStatsResponse, error)
	GetShopInsight(req dto.GetShopInsightRequest) (*dto.GetShopInsightResponse, error)
	GetShopProfile(userId int) (*dto.ShopProfile, error)
	UpdateShopProfile(userId int, req dto.ShopProfile) error
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
	previousShop, err := s.FindShopByUserId(userID)
	if err != nil && !errors.Is(err, commonError.ErrShopNotFound) {
		return nil, err
	}

	if previousShop != nil {
		return nil, commonError.ErrUserHasShop
	}

	request.Name = strings.Trim(request.Name, " ")

	if isShopNameValid := shopUtils.ValidateShopName(request.Name); !isShopNameValid {
		return nil, commonError.ErrInvalidShopName
	}

	courierServices, err := s.courierServiceService.GetCourierServicesByCourierIDs(request.CourierIDs)
	if err != nil {
		return nil, err
	}

	shop := model.Shop{
		Name:           request.Name,
		PhotoUrl:       request.PhotoUrl,
		JoinedDate:     time.Now(),
		Slug:           stringUtils.GenerateSlug(request.Name),
		CourierService: courierServices,
		AddressID:      request.AddressID,
		UserID:         userID,
	}

	err = s.shopRepository.Create(&shop)
	if err != nil {
		return nil, err
	}

	return &shop, nil
}

func (s *shopServiceImpl) GetShopStats(userId int) (*dto.GetShopStatsResponse, error) {
	shop, err := s.shopRepository.FindShopByUserId(userId)
	if err != nil {
		return nil, err
	}

	return s.shopRepository.GetShopStats(shop.ID)
}

func (s *shopServiceImpl) GetShopInsight(req dto.GetShopInsightRequest) (*dto.GetShopInsightResponse, error) {
	shop, err := s.shopRepository.FindShopByUserId(req.UserId)
	if err != nil {
		return nil, err
	}

	return s.shopRepository.GetShopInsight(shop.ID, req)
}

func (s *shopServiceImpl) GetShopProfile(userId int) (*dto.ShopProfile, error) {
	shop, err := s.shopRepository.FindShopByUserId(userId)
	if err != nil {
		return nil, err
	}

	return dto.ComposeShopProfileFromModel(shop), nil
}

func (s *shopServiceImpl) UpdateShopProfile(userId int, req dto.ShopProfile) error {
	shop, err := s.shopRepository.FindShopByUserId(userId)
	if err != nil {
		return err
	}

	req.ComposeToModel(shop)
	return s.shopRepository.UpdateShop(shop)
}
