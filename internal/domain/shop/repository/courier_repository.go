package repository

import (
	"errors"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type CourierRepository interface {
	GetByShopID(shopID int) ([]*model.Courier, error)
	GetByServiceIDAndShopID(courierID, shopID int) (*model.Courier, error)
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

func (r *courierRepositoryImpl) GetByServiceIDAndShopID(courierServiceID, shopID int) (*model.Courier, error) {
	var courier model.Courier

	err := r.db.
		Joins("JOIN courier_services cs ON couriers.id = cs.courier_id").
		Joins("JOIN shop_couriers ON cs.id = shop_couriers.courier_service_id").
		Where("shop_couriers.shop_id = ?", shopID).
		Where("cs.id = ?", courierServiceID).
		First(&courier).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commonErr.ErrCourierNotFound
		}

		return nil, err
	}

	return &courier, nil
}
