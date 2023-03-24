package repository

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type CourierServiceRepository interface {
	GetByCourierIDs(courierIDs []int) ([]*model.CourierService, error)
	CreateCourierService(tx *gorm.DB ,courierService []*model.CourierService) error
}

type courierServiceRepositoryImpl struct {
	db *gorm.DB
}

type CourierServiceRConfig struct {
	DB *gorm.DB
}

func NewCourierServiceRepository(cfg *CourierServiceRConfig) CourierServiceRepository {
	return &courierServiceRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *courierServiceRepositoryImpl) GetByCourierIDs(courierIDs []int) ([]*model.CourierService, error) {
	var courierServices []*model.CourierService

	err := r.db.Where("courier_id IN ?", courierIDs).Find(&courierServices).Error
	if err != nil {
		return nil, err
	}

	return courierServices, nil
}


func (r *courierServiceRepositoryImpl) CreateCourierService(tx *gorm.DB, courierService []*model.CourierService) error {
	err := tx.Create(&courierService).Error
	if err != nil {
		return err
	}

	return nil
}