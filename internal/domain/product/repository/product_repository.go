package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetByCodeFull(code string) (*model.Product, error)
}

type productRepositoryImpl struct {
	db *gorm.DB
}

type ProductRConfig struct {
	DB *gorm.DB
}

func NewProductRepository(cfg *ProductRConfig) ProductRepository {
	return &productRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *productRepositoryImpl) GetByCodeFull(code string) (*model.Product, error) {
	var product model.Product

	err := r.db.Where("code = ?", code).Preload("BulkPrice").Preload("VariantGroup").Preload("Media").First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrProductDoesNotExist
		}

		return nil, err
	}

	return &product, nil
}
