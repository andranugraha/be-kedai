package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
)

type ShopPromotionService interface {
	GetSellerPromotions(userID int, request *dto.SellerPromotionFilterRequest) (*commonDto.PaginationResponse, error)
	GetSellerPromotionById(userId int, promotionId int) (*dto.SellerPromotion, error)
}

type shopPromotionServiceImpl struct {
	shopPromotionRepository repository.ShopPromotionRepository
	shopService             ShopService
}

type ShopPromotionSConfig struct {
	ShopPromotionRepository repository.ShopPromotionRepository
	ShopService             ShopService
}

func NewShopPromotionService(cfg *ShopPromotionSConfig) ShopPromotionService {
	return &shopPromotionServiceImpl{
		shopPromotionRepository: cfg.ShopPromotionRepository,
		shopService:             cfg.ShopService,
	}
}

func (s *shopPromotionServiceImpl) GetSellerPromotions(userID int, request *dto.SellerPromotionFilterRequest) (*commonDto.PaginationResponse, error) {
	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return nil, err
	}

	promotions, totalRows, totalPages, err := s.shopPromotionRepository.GetSellerPromotions(shop.ID, request)
	if err != nil {
		return nil, err
	}

	return &commonDto.PaginationResponse{
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Page:       request.Page,
		Limit:      request.Limit,
		Data:       promotions,
	}, nil
}

func (s *shopPromotionServiceImpl) GetSellerPromotionById(userId int, promotionId int) (*dto.SellerPromotion, error) {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return nil, err
	}

	promotion, err := s.shopPromotionRepository.GetSellerPromotionById(shop.ID, promotionId)
	if err != nil {
		return nil, err
	}

	return promotion, nil
}
