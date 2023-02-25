package repository

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type UserAddressRepository interface {
	AddUserAddress(*model.UserAddress) (*model.UserAddress, error)
	GetAllUserAddress(userId int) ([]*model.UserAddress, error)
}

type userAddressRepository struct {
	db *gorm.DB
}

type UserAddressRConfig struct {
	DB *gorm.DB
}

func NewUserAddressRepository(cfg *UserAddressRConfig) UserAddressRepository {
	return &userAddressRepository{
		db: cfg.DB,
	}
}

func (r *userAddressRepository) AddUserAddress(newAddress *model.UserAddress) (*model.UserAddress, error) {
	var totalRows int64
	var maxAddress int64 = 10

	err := r.db.Model(&model.UserAddress{}).Where("user_id = ?", newAddress.UserID).Count(&totalRows).Error
	if err != nil {
		return nil, err
	}

	if totalRows >= maxAddress {
		return nil, errs.ErrMaxAddress
	}

	err = r.db.Create(newAddress).Error
	if err != nil {
		return nil, err
	}

	return newAddress, nil
}

func (r *userAddressRepository) GetAllUserAddress(userId int) ([]*model.UserAddress, error) {
	var addresses []*model.UserAddress

	err := r.db.Where("user_id = ?", userId).
		Preload("Subdistrict").
		Preload("District").
		Preload("City").
		Preload("Province").
		Find(&addresses).Error
	if err != nil {
		return nil, err
	}

	return addresses, nil
}
