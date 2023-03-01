package repository

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type CourierRepository interface {
	GetByShopID(shopID int) ([]*model.Courier, error)
}

type courierRepositoryImpl struct {
	db *gorm.DB
}

type CourierRConfig struct {
	DB *gorm.DB
}

func NewCourierRepository(cfg *CourierRConfig) CourierRepository {
	return &courierRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *courierRepositoryImpl) GetByShopID(shopID int) ([]*model.Courier, error) {
	var couriers []*model.Courier

	err := r.db.
		Joins("join courier_services cs ON couriers.id = cs.id").
		Joins("right join shop_couriers ON cs.id = shop_couriers.id AND shop_couriers.shop_id = ?", shopID).
		Find(&couriers).
		Error

	if err != nil {
		return nil, err
	}

	return couriers, nil
}
