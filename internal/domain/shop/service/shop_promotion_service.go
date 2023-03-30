package service

import (
	"kedai/backend/be-kedai/internal/common/constant"
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
	CreateShopPromotion(userID int, request *dto.CreateShopPromotionRequest) (*dto.CreateShopPromotionResponse, error)
	DeletePromotion(userID int, promotionID int) error
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

	promotion, err := s.shopPromotionRepository.GetSellerPromotionById(shop.ID, promotionId)
	if err != nil {
		return err
	}

	var shopPromotion *model.ShopPromotion
	if promotion.Status == constant.VoucherPromotionStatusOngoing {
		if !req.StartPeriod.IsZero() {
			shopPromotion = &model.ShopPromotion{
				ID:          promotion.ID,
				Name:        req.Name,
				StartPeriod: promotion.StartPeriod,
				EndPeriod:   req.EndPeriod,
				ShopId:      shop.ID,
			}
		}
	} else {
		shopPromotion = &model.ShopPromotion{
			ID:          promotion.ID,
			Name:        req.Name,
			StartPeriod: req.StartPeriod,
			EndPeriod:   req.EndPeriod,
			ShopId:      shop.ID,
		}
	}

	if shopPromotion.Name == "" {
		shopPromotion.Name = promotion.Name
	}
	if err := req.ValidateDateRange(promotion.StartPeriod, promotion.EndPeriod); err != nil {
		return err
	}

	var productPromotions []*productModel.ProductPromotion
	for _, products := range promotion.Product {
		for _, skus := range products.SKUs {
			productPromotionID := skus.Promotion.ID

			for _, pp := range req.ProductPromotions {
				if pp.Type == "" {
					pp.Type = skus.Promotion.Type
				}
				if pp.Amount == 0 {
					pp.Amount = skus.Promotion.Amount
				}
				if pp.Stock == 0 {
					pp.Stock = skus.Promotion.Stock
				}
				if pp.PurchaseLimit == 0 {
					pp.PurchaseLimit = skus.Promotion.PurchaseLimit
				}
				if pp.SkuId == 0 {
					pp.SkuId = skus.Promotion.SkuId
				}

				productPromotions = append(productPromotions, &productModel.ProductPromotion{
					ID:            productPromotionID,
					Type:          pp.Type,
					Amount:        pp.Amount,
					Stock:         pp.Stock,
					PurchaseLimit: pp.PurchaseLimit,
					SkuId:         pp.SkuId,
					PromotionId:   promotion.ID,
				})
			}
		}
	}

	return s.shopPromotionRepository.Update(shopPromotion, productPromotions)
}

func (s *shopPromotionServiceImpl) CreateShopPromotion(userID int, request *dto.CreateShopPromotionRequest) (*dto.CreateShopPromotionResponse, error) {
	if isPromotionNameValid := productUtils.ValidateProductName(request.Name); !isPromotionNameValid {
		return nil, errs.ErrInvalidPromotionNamePattern
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

func (s *shopPromotionServiceImpl) DeletePromotion(userID int, promotionID int) error {
	shop, err := s.shopService.FindShopByUserId(userID)
	if err != nil {
		return err
	}

	err = s.shopPromotionRepository.Delete(shop.ID, promotionID)
	if err != nil {
		return err
	}

	return nil
}
