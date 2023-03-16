package repository

import (
	"kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
)

type ProductVariantRepository interface {
	Create(tx *gorm.DB, payload []*model.ProductVariant) error
}

type productVariantRepositoryImpl struct {
	db *gorm.DB
}

type ProductVariantRConfig struct {
	DB *gorm.DB
}

func NewProductVariantRepository(cfg *ProductVariantRConfig) ProductVariantRepository {
	return &productVariantRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *productVariantRepositoryImpl) Create(tx *gorm.DB, payload []*model.ProductVariant) error {
	return tx.Create(&payload).Error
}
