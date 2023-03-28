package repository

import (
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	productRepo "kedai/backend/be-kedai/internal/domain/product/repository"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type ShopPromotionRepository interface {
	Create(shopID int, request *dto.CreateShopPromotionRequest) (*model.ShopPromotion, error)
}

type shopPromotionRepositoryImpl struct {
	db                *gorm.DB
	productRepository productRepo.ProductRepository
}

type ShopPromotionRConfig struct {
	DB                *gorm.DB
	ProductRepository productRepo.ProductRepository
}

func NewShopPromotionRepository(cfg *ShopPromotionRConfig) ShopPromotionRepository {
	return &shopPromotionRepositoryImpl{
		db:                cfg.DB,
		productRepository: cfg.ProductRepository,
	}
}

func (r *shopPromotionRepositoryImpl) Create(shopID int, request *dto.CreateShopPromotionRequest) (*model.ShopPromotion, error) {
	tx := r.db.Begin()
	defer tx.Commit()

	shopPromotion := request.GenerateShopPromotion()
	shopPromotion.ShopId = shopID

	err := tx.Create(shopPromotion).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	productPromotions := []*productModel.ProductPromotion{}

	for _, pp := range request.ProductPromotions {
		productPromotions = append(productPromotions, &productModel.ProductPromotion{
			Type:          pp.Type,
			Amount:        pp.Amount,
			Stock:         pp.Stock,
			IsActive:      *pp.IsActive,
			PurchaseLimit: pp.PurchaseLimit,
			SkuId:         pp.SkuId,
			PromotionId:   shopPromotion.ID,
		})
	}

	err = tx.Create(productPromotions).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return shopPromotion, nil
}
