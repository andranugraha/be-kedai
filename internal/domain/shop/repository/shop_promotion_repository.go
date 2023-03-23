package repository

import "gorm.io/gorm"

type ShopPromotionRepository interface {
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

func (r *shopPromotionRepositoryImpl) GetSellerPromotions() {

}
