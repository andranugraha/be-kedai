package repository

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type CourierRepository interface {
	GetByShopID(shopID int) ([]*model.Courier, error)
	GetByProductID(productID int) ([]*model.Courier, error)
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
		Joins("JOIN courier_services cs ON couriers.id = cs.courier_id").
		Joins("JOIN shop_couriers ON cs.id = shop_couriers.courier_service_id").
		Where("shop_couriers.shop_id = ?", shopID).
		Distinct().
		Find(&couriers).
		Error

	if err != nil {
		return nil, err
	}

	return couriers, nil
}

func (r *courierRepositoryImpl) GetByProductID(productID int) ([]*model.Courier, error) {
	var couriers []*model.Courier

	err := r.db.
		Joins("JOIN courier_services cs ON couriers.id = cs.courier_id").
		Joins("JOIN product_couriers ON cs.id = product_couriers.courier_service_id").
		Where("product_couriers.product_id = ?", productID).
		Distinct().
		Find(&couriers).
		Error

	if err != nil {
		return nil, err
	}

	return couriers, nil
}
