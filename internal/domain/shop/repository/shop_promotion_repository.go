package repository

import (
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/shop/dto"

	"gorm.io/gorm"
)

type ShopPromotionRepository interface {
	Create(shopID int, request *dto.CreateShopPromotionRequest) (*dto.CreateShopPromotionResponse, error)
}

type shopPromotionRepositoryImpl struct {
	db *gorm.DB
}

type ShopPromotionRConfig struct {
	DB *gorm.DB
}

func NewShopPromotionRepository(cfg *ShopPromotionRConfig) ShopPromotionRepository {
	return &shopPromotionRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *shopPromotionRepositoryImpl) Create(shopID int, request *dto.CreateShopPromotionRequest) (*dto.CreateShopPromotionResponse, error) {
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

	response := &dto.CreateShopPromotionResponse{
		ShopPromotion:     *shopPromotion,
		ProductPromotions: productPromotions,
	}

	return response, nil
}
