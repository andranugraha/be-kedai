package repository

import (
	"errors"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type CourierRepository interface {
	GetShipmentList(shopId int) ([]*dto.ShipmentCourierResponse, error)
	GetByShopID(shopID int) ([]*model.Courier, error)
	GetByServiceIDAndShopID(courierID, shopID int) (*model.Courier, error)
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

func (r *courierRepositoryImpl) GetShipmentList(shopId int) ([]*dto.ShipmentCourierResponse, error) {
	var couriers []*dto.ShipmentCourierResponse

	db := r.db.Select(`DISTINCT ON (couriers.id) couriers.*, COALESCE(sc.is_active, false) as is_active`).
		Joins("JOIN courier_services cs ON cs.courier_id = couriers.id").
		Joins("LEFT JOIN shop_couriers sc ON sc.courier_service_id = cs.id AND sc.shop_id = ?", shopId)

	err := db.Model(&model.Courier{}).Find(&couriers).Error
	if err != nil {
		return nil, err
	}

	return couriers, nil
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
