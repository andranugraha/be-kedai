package repository

import (
	model "kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetByCode(Code string) (*model.Product, error)
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

func (r *productRepositoryImpl) GetByCode(Code string) (*model.Product, error) {
	var product model.Product

	err := r.db.Where("code = ?", Code).First(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}
