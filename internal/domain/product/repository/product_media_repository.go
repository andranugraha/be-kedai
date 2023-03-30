package repository

import (
	"kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
)

type ProductMediaRepository interface {
	Delete(tx *gorm.DB, productId int) error
}

type productMediaRepositoryImpl struct {
	db *gorm.DB
}

type ProductMediaRConfig struct {
	DB *gorm.DB
}

func NewProductMediaRepository(cfg *ProductMediaRConfig) ProductMediaRepository {
	return &productMediaRepositoryImpl{
		db: cfg.DB,
	}
}

func (p *productMediaRepositoryImpl) Delete(tx *gorm.DB, productId int) error {

	err := tx.Model(&model.ProductMedia{}).Where("product_id = ?", productId).Unscoped().Delete(&model.ProductMedia{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
