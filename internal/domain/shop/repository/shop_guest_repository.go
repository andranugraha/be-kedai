package repository

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"

	commonErr "kedai/backend/be-kedai/internal/common/error"

	"gorm.io/gorm"
)

type ShopGuestRepository interface {
	CreateShopGuest(shopGuest *model.ShopGuest) (*model.ShopGuest, error)
}

type shopGuestRepositoryImpl struct {
	db *gorm.DB
}

type ShopGuestRConfig struct {
	DB *gorm.DB
}

func NewShopGuestRepository(cfg *ShopGuestRConfig) ShopGuestRepository {
	return &shopGuestRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *shopGuestRepositoryImpl) CreateShopGuest(shopGuest *model.ShopGuest) (*model.ShopGuest, error) {
	err := r.db.Create(shopGuest).Error
	if err != nil {
		if commonErr.IsForeignKeyError(err) {
			return nil, commonErr.ErrShopNotFound
		}
		return nil, err
	}

	return shopGuest, nil
}
