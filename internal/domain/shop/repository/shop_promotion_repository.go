package repository

import (
	productRepo "kedai/backend/be-kedai/internal/domain/product/repository"

	"gorm.io/gorm"
)

type ShopPromotionRepository interface {
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
