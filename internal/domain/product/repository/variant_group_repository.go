package repository

import (
	"kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
)

type VariantGroupRepository interface {
	Create(tx *gorm.DB, variantGroups []*model.VariantGroup) error
}

type variantGroupRepositoryImpl struct {
	db *gorm.DB
}

type VariantGroupRConfig struct {
	DB *gorm.DB
}

func NewVariantGroupRepository(cfg *VariantGroupRConfig) VariantGroupRepository {
	return &variantGroupRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *variantGroupRepositoryImpl) Create(tx *gorm.DB, variantGroups []*model.VariantGroup) error {
	return tx.Create(&variantGroups).Error
}
