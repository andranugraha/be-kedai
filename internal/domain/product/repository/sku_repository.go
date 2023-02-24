package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
)

type SkuRepository interface {
	GetByID(ID int) (*model.Sku, error)
}

type skuRepositoryImpl struct {
	db *gorm.DB
}

type SkuRConfig struct {
	DB *gorm.DB
}

func NewSkuRepository(cfg *SkuRConfig) SkuRepository {
	return &skuRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *skuRepositoryImpl) GetByID(ID int) (*model.Sku, error) {
	var sku model.Sku

	err := r.db.Where("id = ?", ID).First(&sku).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrProductDoesNotExist
		}

		return nil, err
	}

	return &sku, err
}
