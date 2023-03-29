package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
	productUtils "kedai/backend/be-kedai/internal/utils/product"
)

type ShopPromotionService interface {
	GetSellerPromotions(userID int, request *dto.SellerPromotionFilterRequest) (*commonDto.PaginationResponse, error)
	GetSellerPromotionById(userId int, promotionId int) (*dto.SellerPromotion, error)
	UpdatePromotion(userId int, promotionId int, req dto.UpdateShopPromotionRequest) error
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

func (s *shopPromotionServiceImpl) UpdatePromotion(userId int, promotionId int, req dto.UpdateShopPromotionRequest) error {
	if isPromotionNameValid := productUtils.ValidateProductName(req.Name); !isPromotionNameValid {
		return errs.ErrInvalidPromotionNamePattern
	}

	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return err
	}

	promotion, err := s.GetSellerPromotionById(shop.ID, promotionId)
	if err != nil {
		return err
	}

	shopPromotion := &model.ShopPromotion{
		ID:          promotion.ID,
		Name:        req.Name,
		StartPeriod: req.StartPeriod,
		EndPeriod:   req.EndPeriod,
		ShopId:      shop.ID,
	}

	var productPromotions []*productModel.ProductPromotion
	for _, products := range promotion.Product {
		for _, skus := range products.SKUs {
			productPromotionID := skus.Promotion.ID

			for _, pp := range req.ProductPromotions {
				productPromotions = append(productPromotions, &productModel.ProductPromotion{
					ID:            productPromotionID,
					Type:          pp.Type,
					Amount:        pp.Amount,
					Stock:         pp.Stock,
					PurchaseLimit: pp.PurchaseLimit,
					SkuId:         pp.PurchaseLimit,
					PromotionId:   promotion.ID,
				})
			}
		}
	}

	return s.shopPromotionRepository.Update(shopPromotion, productPromotions)
}
