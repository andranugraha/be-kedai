package repository

import (
	"kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type UserAddressRepository interface {
	AddUserAddress(*model.UserAddress) (*model.UserAddress, error)
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
	err := r.db.Create(newAddress).Error
	if err != nil {
		return nil, err
	}

	return newAddress, nil
}
