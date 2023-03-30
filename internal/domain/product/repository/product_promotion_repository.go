package repository

import (
	model "kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
)

type ProductPromotionRepository interface {
	Delete(tx *gorm.DB, skuId int) error
}

type productPromotionRepositoryImpl struct {
	db *gorm.DB
}

type ProductPromotionRConfig struct {
	DB *gorm.DB
}

func NewProductPromotionRepository(cfg *ProductPromotionRConfig) ProductPromotionRepository {
	return &productPromotionRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *productPromotionRepositoryImpl) Delete(tx *gorm.DB, skuId int) error {

	if err := tx.Model(&model.ProductPromotion{}).
		Where("sku_id = ?", skuId).
		Delete(&model.ProductPromotion{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
