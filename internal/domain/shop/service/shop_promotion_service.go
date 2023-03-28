package service

import (
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
	productUtils "kedai/backend/be-kedai/internal/utils/product"
)

type ShopPromotionService interface {
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

func (s *shopPromotionServiceImpl) CreateShopPromotion(userID int, request *dto.CreateShopPromotionRequest) (*model.ShopPromotion, error) {
	if isPromotionNameValid := productUtils.ValidateProductName(request.Name); !isPromotionNameValid {
		return nil, commonErr.ErrInvalidPromotionNamePattern
	}

	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return nil, err
	}

	promotion, err := s.shopPromotionRepository.Create(shop.ID, request)
	if err != nil {
		return nil, err
	}

	return promotion, nil
}
